package server

import (
    "chat/client"
    "chat/types"
    "errors"
    "fmt"
    "net"
    "os"
    "strconv"
    "strings"
    "sync/atomic"
    "time"
)

// func New {{{

// Initializes and returns a new Server.
func New(ip string, port int) *Server {
    var s Server
    var err error
    var ids []uint32

    // Set our ids variable for the server
    s.ids = ids

    // Set our server TCP Address
    s.bindy = net.TCPAddr{
        IP: net.ParseIP(ip),
        Port: port,
    }

    // Create a new TCP Listener using our tcpAddr
    if s.listener, err = net.ListenTCP("tcp", &s.bindy); err != nil {
        fmt.Printf("server.New: error - net.ListenTCP(%s): %v", s.bindy, err)
        os.Exit(-1)
    }

    // Return the new server
    return &s
} // }}}

// func s.SetApplication {{{
//
// Sets the application of the server, since we make the server
// prior to making the application
func (s *Server) SetApplication(app types.Application) {
    s.app = app
} // }}}

// func s.Listen {{{
//
// Uses the servers listener to continuously accept incoming TCP connections
func (s *Server) Listen() {
    var errs int
    for {
        conn, err := s.listener.AcceptTCP()

        // Any kind of error?
        if err != nil {
            // Closed?
			if errors.Is(err, net.ErrClosed) {
				// This basically means we are shutting down, so not a real "error".
				return
			}

            // If we get too many errors in a row that we are not programmed
            // here to handle, we close the listening socket and return.
            //
			// The error count is reset when we get a new connection.
			if errs > 5 {
				s.app.OutErr("Too many accept errors, unable to accept any new connections\nPlease enter a command: ")

                s.mu.Lock()
				s.listener.Close()
                s.mu.Unlock()

				return
			}

            // We don't know what this error is, but its not a closed socket?
			//
			// Log it and attempt to continue.
			s.app.OutErr("AcceptTCP(%s): %s\nPlease enter a command: ", s.listener.Addr(), err)

            // Increase the error count, since we do not specifically
            // handle this error, we are unsure if its safe to continue
            // listening or not.
            errs++
        }

        // We have a new connection with no errors, so reset the error counter if needed.
		errs = 0

        // Let's get the connections address
        remoteAddr := conn.RemoteAddr()
        connAddr := strings.Split(remoteAddr.String(), ":")

        // Before we make a new Client and it to our sync map, let's see if this
        // connection already exists for whatever reason
        if s.checkExisting(connAddr[0], connAddr[1]) {
            s.app.OutErr("s.Listen: refusing connection from %s:%s! connection already exists!\nPlease enter a command: ", connAddr[0], connAddr[1])
            continue
        }

        // Lets get the id for our new connection, using atmoics!
        id := atomic.AddUint32(&s.nextID, 1)

        // Createa a new client
        c := client.New(conn, connAddr, id, s.app)

        // Add the new client to our sync map
        s.conns.Store(c.Id, c)

        // Add the new ID to our ids array -- this is only used so
        // we can get a sorted list of connections when List() is called
        s.mu.Lock()
        s.ids = append(s.ids, c.Id)
        s.mu.Unlock()

        // Inform the user of the new connection
        s.app.Out("\nNew incoming connection: %v | %v:%v\n\nPlease enter a command: ", c.Id, c.IP, c.Port)

        // Handle the connection in a new goroutine.
        //
        // The loop then returns to accepting, so that
        // multiple connections may be served concurrently.
        go func() {
            // Ensure we close our connection and remove it from the sync
            // map once HandleClient returns
            defer func() {
                // The only time we should return is when the connection
    			// is closed, or some otherwise unrecoverable error.
    			//
    			// So we remove ourself from the connect list
                // automatically when we return here.
    			s.conns.Delete(c.Id)

                // We also (just in case), call Close() on the connection
                // again,  as this will handle any other errors that don't
    			// end up closing the connection properly.
    			//
    			// Note this is safe to call on an already closed connection.
    			c.Conn.Close()
            }()
            c.HandleClient()
        }()
    }
} // }}}

// func s.checkExisting {{{
//
// Checks if the the connection attempting to be establed
// already exists or not
func (s *Server) checkExisting(ip, port string) bool {
    found := false

    // Range over our connections sync map so we can check
    // the conns ip and port values
    //
    // This provides us with a k, v pair that are interfaces
    // so we need to type assert our value to *client.Client,
    // it's actual data type.
    s.conns.Range(func(_, v interface {}) bool {
        // Type assert the loaded value to the correct type
        c, ok := v.(*client.Client)
        if !ok {
            s.app.OutErr("s.checkExisting: error asserting client type\nPlease enter a command: ")
            return false
        }

        // Are the IP and port the same as the ones provided?
        if c.IP == ip && c.Port == port {
            found = true
            return found
        }
        return found
    })
    return found
} // }}}

// func s.Connect {{{
//
// Attempts to establish a new connection, returning an error should
// anything go wrong
func (s *Server) Connect(destination, port string) error {
    connErr := errors.New("s.Connect: connection already exists!")
    selfErr := errors.New("s.Connect: self connections not allowed!")
    invIpErr := errors.New(fmt.Sprintf("s.Connect: invalid ip given! %s:%s", destination, port))
    invPortErr := errors.New(fmt.Sprintf("s.Connect: invalid port given! Dial request timed out %s:%s", destination, port))

    // Does this connection already exist?
    if s.checkExisting(destination, port) {
        return connErr
    }

    // Are we trying to establish a connnection on our own port?
    p, _ := strconv.ParseInt(port, 10, 64)
    if s.bindy.Port == int(p) {
        return selfErr
    }

    // Were we given an invalid IP address?
    ip := net.ParseIP(destination)
    if ip == nil {
        return invIpErr
    }

    // Create a new net Dialer and set the timeout to be 10 seconds
    // Timeout is max time allowed to wait for a dial to connect
    // (was 20s but that felt painfully slow)
    //
    // We're using a timeout so we don't completely break the program
    // if we never get a new connection cos the user didn't give us a
    // valid IP
    timeout, _ := time.ParseDuration("10s")
    dialer := net.Dialer{ Timeout: timeout }

    // Dial the connection adddress to establish connection.
    tcpAddr := net.JoinHostPort(destination, port)
    conn, err := dialer.Dial("tcp", tcpAddr)
    if err != nil {
        // We timed out, most likey due to an inavlid IP/port combo
        return invPortErr
    }

    // We went a net.TCPConn, not net.Conn so we've gotta do some type assertion
    tcpConn, ok := conn.(*net.TCPConn)
    if !ok {
        return errors.New("s.Connect: error asserting connection type")
    }

    // Lets get the id for our new connection
    id := atomic.AddUint32(&s.nextID, 1)

    // Let's get the connections address
    remoteAddr := conn.RemoteAddr()
    connAddr := strings.Split(remoteAddr.String(), ":")

    // Createa a new client
    c := client.New(tcpConn, connAddr, id, s.app)

    // Add the new client to our sync map
    s.conns.Store(c.Id, c)

    // Add the new ID to our ids array -- this is only used so
    // we can get a sorted list of connections when List() is called
    s.mu.Lock()
    s.ids = append(s.ids, c.Id)
    s.mu.Unlock()

    // Inform the user of the new connection
    s.app.Out("\nNew connection established: %v:%v\n", c.IP, c.Port)

    // Handle the connection in a new goroutine.
    //
    // The loop then returns to accepting, so that
    // multiple connections may be served concurrently.
    go func() {
        // Ensure we close our connection and remove it from the sync
        // map once HandleClient returns
        defer func() {
            // The only time we should return is when the connection
            // is closed, or some otherwise unrecoverable error.
            //
            // So we remove ourself from the connect list
            // automatically when we return here.
            s.conns.Delete(c.Id)

            // We also (just in case), call Close() on the connection
            // again,  as this will handle any other errors that don't
            // end up closing the connection properly.
            //
            // Note this is safe to call on an already closed connection.
            c.Conn.Close()
        }()
        c.HandleClient()
    }()

    return nil
} // }}}

// func s.List {{{
//
// Lists the IP addresses and port numbers associated with all
// currently established connections
func (s *Server) List() {
    s.app.Out("id |  IP Address   | Port\n")
    s.app.Out("---+---------------+-----\n")

    // Let's grab the ids the server currently has
    s.mu.Lock()
    ids := s.ids
    s.mu.Unlock()

    // Range over our arry of ids and print them.
    for _, id := range ids {
        // Try loading our connection from our sync map
        v, ok := s.conns.Load(id)

        // Check if it loaded or not - if it didn't its likely been
        // deleted from the map so just continue
        if !ok {
            continue
        }

        // Type assert the loaded value to the correct type
        c, ok := v.(*client.Client)
        if !ok {
            s.app.OutErr("s.List: error asserting client type\nPlease enter a command: ")
            return
        }

        // Print the connection details
        s.app.Out(" %d | %s | %s\n", c.Id, c.IP, c.Port)
    }
} // }}}

// func s.Terminate {{{
//
// Terminates the connection associated with the given connection id,
// returning an error should anything go wrong
func (s *Server) Terminate(conn int) error {
    invInput := fmt.Sprintf("s.Terminate: must give a valid connection ID! Use list to see a list of all current connections.")
    invInputErr := errors.New(invInput)

    // Try loading our connection from our sync map
    v, ok := s.conns.Load(uint32(conn))

    // Were we given a valid connection ID?
    if !ok {
        return invInputErr
    }

    // Type assert the loaded value to the correct type
    c, ok := v.(*client.Client)
    if !ok {
        return errors.New("s.Terminate: error asserting client type")
    }

    // Attempt to close the connection
    //
    // We don't need to do anything to remove the connection from the
    // connections map because it will already be removed once the client
    // returns from the HandleClient func
    if err := c.CloseConn(); err != nil {
        return err
    }

    return nil
} // }}}

// func s.Send {{{
//
// Attempts to send a given message to the connection associated with the
// given connection id, returning an error should anything go wrong
func (s *Server) Send(conn int, message string) error {
    invInput := fmt.Sprintf("s.Send: %d invalid ID! Use list to see a list of all current connections", conn)
    invInputErr := errors.New(invInput)

    tooLong := fmt.Sprintf("s.Send: message is too long! Max length a message can be is 100 characters. Your message is %d characters", len(message))
    tooLongErr := errors.New(tooLong)

    // Were we given a message of valid length?
    if len(message) > 100 {
        return tooLongErr
    }

    // Load our connection from the provided connection ID
    v, ok := s.conns.Load(uint32(conn))

    // Were we given a valid connection ID?
    if !ok {
        return invInputErr
    }

    // Type assert the loaded value to the correct type
    c, ok := v.(*client.Client)
    if !ok {
        return errors.New("s.Send: error asserting client type")
    }

    // Lets create a byte array of our message to send to the connection
    msg := []byte(message)

    // Okay now that we've loaded the connection, let's try to send
    // the message
    c.Conn.Write(msg)
    s.app.Out("Message sent to connection %d!\n", c.Id)
    return nil
} // }}}

// func s.Exit {{{
//
// Exits the program, closing any established connections prior to doing so
func (s *Server) Exit() {
    s.app.Out("Closing any established connections .. \n")

    s.conns.Range(func(k, v interface {}) bool {
        // Type assert the loaded value to the correct type
        c, ok := v.(*client.Client)
        if !ok {
            s.app.OutErr("s.Exit: error asserting client type\n")
            return false
        }

        // Try and close the connection
        if err := c.CloseConn(); err != nil {
            s.app.OutErr("s.Exit: error closing connection %d! err:=%v\n", c.Id, err)
            return false
        }
        return true
    })

    // Stop listening for new connections ..
    s.mu.Lock()
    s.listener.Close()
    s.mu.Unlock()

    // Let the user know we're shutting down now and exit
    s.app.Out("Exiting program now .. bye!\n")
    os.Exit(0)
}  // }}}

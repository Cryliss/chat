// Package app provides user input functionality
package app

import (
    "github.com/Cryliss/chat/types"
    "errors"
    "fmt"
    "net"
    "os"
    "strings"
    "strconv"
)

// func New {{{

// New Initializes & returns a new application and any errors that may have occurred.
func New(port int, ip string, server types.Server) (*Application, error) {
    // Let's generate our service string, in the format "{ip}:{port}"
    portStr := fmt.Sprintf("%d", port)
    service := net.JoinHostPort(ip, portStr)

    a := Application{
        s: server,
        port: port,
        ip: ip,
        service: service,
    }
    a.startupText()

    return &a, nil
} // }}}

// func a.Out {{{

// Out Prints message to the standard output device
func (a *Application) Out(format string, b ...interface{}) {
    // We're we given any variables that should be added to the string?
    if b == nil {
        // No? Okay, let's not add them to Fprint, otherwise we get errors :D
        fmt.Fprintf(os.Stdout, format)
        return
    }
    fmt.Fprintf(os.Stdout, format, b...)
} // }}}

// func a.OutErr {{{
//
// OutErr Prints message to the standard output error device
func (a *Application) OutErr(format string, b ...interface{}) {
    // We're we given any variables that should be added to the string?
    if b == nil {
        // No? Okay, let's not add them to Fprint, otherwise we get errors :D
        fmt.Fprintf(os.Stderr, format)
        return
    }

    fmt.Fprintf(os.Stderr, format, b...)
} // }}}

// func a.startupText {{{

// Prints the text that should be displayed on application startup
func (a *Application) startupText() {
    sText := `
CHATTY: A Chat Application for Remote Message Exchange
------------------------------------------------------
Available commands:
    1. help
    2. myip
    3. myport
    4. connect <destination> <port no>
    5. list
    6. terminate <connection id>
    7. send <connection id> <message>
    8. exit

You may either type the command name, i.e. 'connect <destination> <port no>', or the command number, i.e. '4 <destination> <port no>'
Type 'help' for an explanation of each command, or type 'help <command>' to get the explanation for a specific command
`
    a.Out("%s", sText)
} // }}}

// func a.ParseInput {{{

// ParseInput Parses the users input and calls the function associated with
// the given command
func (a *Application) ParseInput(userInput string) error {
    // Split the users input into array of strings
    inputArgs := strings.SplitN(userInput, " ", 3)
    numArgs := len(inputArgs)

    // Create some common errors for input mistakes we may see
    connErr := errors.New("connect input error: You must give both the destination and the port when using the connect command")
    termErr := errors.New("terminate input error: You must give the connection id you wish to terminate\nType `list` to get a list of connections and their ids")
    sendErr := errors.New("send input error: You must give both the connection id and a message to the connection")
    inputErr := errors.New("invalid input error: You must give one of the accepted app commands\nType 'help' to get a list of available commands")

    // Check what the command was, the first item in the input, and
    // perform the actions necessary for that command
    switch inputArgs[0] {
    case "1":
        fallthrough
    case "help":
        // If we have more than 2 input args, that means they want to see
        // the command information for multiple commands, so let's loop
        // through each of them
        if numArgs > 2 {
            for i := 1; i < numArgs; i++ {
                a.help(inputArgs[i])
            }
            return nil
        }

        // If we only have 2, then the user just wants to see this one command
        if numArgs == 2 {
            a.help(inputArgs[1])
            return nil
        }

        // We only had help in our input, so they want the full list
        a.help("")
        return nil
    case "2":
        fallthrough
    case "myip":
        a.myip()
        return nil
    case "3":
        fallthrough
    case "myport":
        a.myport()
        return nil
    case "4":
        fallthrough
    case "connect":
        // Do we have the proper number of arguments?
        if numArgs < 3 {
            return connErr
        }

        // Yes, so let's grab them from the array and attempt
        // to establish thew new connection
        destination := inputArgs[1]
        port := inputArgs[2]
        if err := a.s.Connect(destination, port); err != nil {
            return err
        }
        return nil
    case "5":
        fallthrough
    case "list":
        a.s.List()
        return nil
    case "6":
        fallthrough
    case "terminate":
        // Do we have the proper number of arguments?
        if numArgs != 2 {
            return termErr
        }

        // Yes, so let's get the connection that we need to terminate
        // and attempt to terminate it
        conn, _ := strconv.ParseInt(inputArgs[1], 10, 64)
        if err := a.s.Terminate(int(conn)); err != nil {
            return err
        }
        return nil
    case "7":
        fallthrough
    case "send":
        // Do we have the proper number of arguments?
        if numArgs != 3 {
            return sendErr
        }

        // Yes, so let's get the connection we need to send the message to,
        // as well as the message itself
        conn, _ := strconv.ParseInt(inputArgs[1], 10, 64)
        msg := inputArgs[2]

        // Now let's attempt to send the message
        if err := a.s.Send(int(conn), msg); err != nil {
            return err
        }
        return nil
    case "8":
        fallthrough
    case "exit":
        a.s.Exit()
        return nil
    default:
        // We didn't find a matching command for their input, let's throw an error
        return inputErr
    }
    return nil
}  // }}}

// func a.help {{{

// help Prints the application commands
func (a *Application) help(command string) {
    a.Out("\nApplication Commands\n")
    a.Out("--------------------\n")

    switch command {
    case "":
        cmds := `1. help - Displays available application commands
2. myip - Displays the IP address of this process
3. myport - Displays the port on which this process is listening for incoming connections
4. connect <destination> <port no> - Establishes a new TCP connection to the specified <destination> at the specified <port no>
5. list - Displays a numbered list of all the connections this process is a part of
   For example:
    id |  IP Address   | Port
    ---+---------------+-----
    1 | 192.168.21.20 | 4545
    2 | 192.168.21.21 | 5454
6. terminate <connection id> - Terminates the connection associated with the given connection id
7. send <connection id> <message> - Sends a message to the host on the connection that is designated by the connection id
8. exit - Closes all connections and terminates the process
`
        a.Out(cmds)
        break
    case "1":
        fallthrough
    case "help":
        a.Out("1. help - Displays available application commands\n")
        break
    case "2":
        fallthrough
    case "myip":
        a.Out("2. myip - Displays the IP address of this process\n")
        break
    case "3":
        fallthrough
    case "myport":
        a.Out("3. myport - Displays the port on which this process is listening for incoming connections\n")
        break
    case "4":
        fallthrough
    case "connect":
        a.Out("4. connect <destination> <port no> - Establishes a new TCP connection to the specified <destination> at the specified <port no>\n")
        break
    case "5":
        fallthrough
    case "list":
        a.Out(`5. list - Displays a numbered list of all the connections this process is a part of
for example:
    id |  IP Address   | Port
    ---+---------------+-----
     1 | 192.168.21.20 | 4545
     2 | 192.168.21.21 | 5454
`)
        break
    case "6":
        fallthrough
    case "terminate":
        a.Out("6. terminate <connection id> - Terminates the connection associated with the given connection id\n")
        break
    case "7":
        fallthrough
    case "send":
        a.Out("7. send <connection id> <message> - Sends a message to the host on the connection that is designated by the connection id\n")
        break
    case "8":
        fallthrough
    case "exit":
        a.Out("8. exit - Closes all connections and terminates the process\n")
        break
    default:
        return
    }
} // }}}

// func a.myip {{{

// myip Prints the users IP address
func (a *Application) myip() {
    a.Out("Your IP address is: %s\n", a.ip)
} // }}}

// func a.myport {{{

// myport Prints the users port number
func (a *Application) myport() {
    a.Out("Your port is: %d\n", a.port)
} // }}}

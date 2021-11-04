package server

import (
    "chat/types"
    "net"
    "sync"
)

// type Server struct {{{

type Server struct {
    // Our application so we can display messages to the user
    //
    // We are using types.Application rather than app.Application
    // because that causes an import loop which is not allowed ...
    app types.Application

    // The tcpAddr we want to bind to & accept connections on
    bindy net.TCPAddr

    // A sync map for our connections
    //
    // This allows us to use atomics to safely get new connection
    // ID values, without fear of data races
    conns sync.Map

    // List of connection IDs .. only using this so that my list
    // of connections will be sorted .. doesn't work that way if I just
    // range over the sync map
    ids []uint32

    // Locks reading on this struct, avoids data races!
    mu sync.Mutex

    // The next availabe connection ID
    //
    // Only access this using atomics!
    nextID uint32

    // Listener that will accept incoming connections
    listener *net.TCPListener
} // }}}

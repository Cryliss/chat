package client

import (
    "chat/types"
    "net"
)

// type Client struct  {{{

type Client struct {
    // Our application so we can print to the user
    app types.Application

    // Assigned ID for the connection
    Id uint32

    // Connections IP Address
    IP string

    // Connections port number
    Port string 

    // The actual connection itself
    Conn *net.TCPConn
} // }}}

// Package client handles new client connections
package client

import (
	"github.com/Cryliss/chat/types"
	"net"
)

// type Client struct  {{{

// Client data type to hold information related to the client connection
type Client struct {
	// Our application so we can print to the user
	app types.Application

	// Assigned ID for the connection
	ID uint32

	// Connections IP Address
	IP string

	// Connections port number
	Port string

	// The actual connection itself
	Conn *net.TCPConn
} // }}}

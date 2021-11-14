package app

import "chat/types"

// Commands the user can give the application
// We are using a map[string]string here, just in case
// the user wants to be lazy and not type the whole thing out
var commands = map[string]string{
    "1": "help",
    "2": "myip",
    "3": "myport",
    "4": "connect",
    "5": "list",
    "6": "terminate",
    "7": "send",
    "8": "exit",
}

type Application struct {
    s           types.Server

    // Map of the application commands, see above
    commands    map[string]string

    // The port provided at runtime
    port        int

    // The ip address of the application
    ip          string

    // To hold the formatted string for the connection service
    // ':<port>'
    service     string
}

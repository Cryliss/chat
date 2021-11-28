package main

import (
    "bufio"
    "github.com/Cryliss/chat/app"
    "github.com/Cryliss/chat/server"
    "flag"
    "fmt"
    "net"
    "os"
    "strings"
)

// func usage {{{

// usage Prints information on how to use the program and then exits
func usage() {
    fmt.Printf("usage: %s\n", os.Args[0])
    flag.PrintDefaults()
    os.Exit(-1)
} // }}}

// func main {{{

func main() {
    var port int

    // Lets load our flags.
    flag.IntVar(&port, "port", -1, "Port to listen on for incoming connections")
    flag.Parse()

    // Did we get a port number?
    if port == -1 {
        usage()
    }

    // Get the ip address of the machine running the program
    p := fmt.Sprintf("%d", port)
    ip := GetOutboundIP(p)

    // Create a new server
    server := server.New(ip, port)

    // Create a new application for user input / output
    app, _ := app.New(port, ip, server)

    // Set the server's app and start listening for connections
    server.SetApplication(app)
    go server.Listen()

    // Create a new bufio reader to read user input from the command line
    reader := bufio.NewReader(os.Stdin)

    // Variable to store our scanned input into
    var userInput string

    for {
        // Prompt the user for a command
        app.Out("\nPlease enter a command: ")

        // Read user input and save into userInput variable
        userInput, _ = reader.ReadString('\n')
        userInput = strings.Replace(userInput, "\n", "", -1)

        // Parse and handle the users input
        // If the request resulted in an error, let's let the user know
        err := app.ParseInput(userInput)
        if err != nil {
            app.OutErr("ERROR %v\n", err)
        }
    }
} // }}}

// func GetOutboundIP() {{{
//
// Get preferred outbound ip of this machine
// src: https://stackoverflow.com/a/37382208
func GetOutboundIP(port string) string {
    s := "8.8.8.8:" + port
    conn, err := net.Dial("udp", s)
    if err != nil {
        fmt.Fprintf(os.Stderr, "%v\n", err)
        os.Exit(-1)
    }
    defer conn.Close()

    localAddr := conn.LocalAddr().(*net.UDPAddr)

    ip := fmt.Sprintf("%v", localAddr.IP)
    ipArr := strings.Split(ip, ":")
    return ipArr[0]
} // }}}

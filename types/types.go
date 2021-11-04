package types

// This package is required to avoid import cycles in Go
//
// It defines the three types of interfaces we have in our progam
// and the functions that are required to be in each of them

type Application interface {
    Out(format string, a ...interface{})
    OutErr(format string, a ...interface{})
    ParseInput(userInput string) error
}

type Client interface {
    HandleClient()
    CloseConn() error
}

type Server interface {
    Connect(destination, port string) error
    List()
    Terminate(conn int) error
    Send(conn int, message string) error
    Exit()
}

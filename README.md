Chatty
------
**COMP 429 Programming Assignment 1**  
A Chat Application for Remote Message Exchange  
*by Sabra Bilodeau*  


## Getting Started with Chatty
### Building the Program
1. [Download Go](https://golang.org/dl/), if you don't already have it installed on your local machine.  
2. Download the source code via one of these two methods  
    a. Clone the repository by running `git clone https://github.com/Cryliss/comp-429.git` in the terminal  
    b. Download the source code directly from GitHub, and place it in your Go source code directory - e.g. `/Users/sabra/go/src`  
3. From terminal, navigate to the chat application directory. e.g. `cd ~/go/src/chat`
4. Run `make build` to build the applications binary file

`make build` builds the program with the `--race` flag in order to help catch any data races that may occur within our application due to the goroutines. (ideally none will, but just in case .. )   
An example (*that has since been corrected!*) is below:

```shell
==================
WARNING: DATA RACE
Write at 0x00c0001247e0 by goroutine 7:
  runtime.mapassign_fast32()
      /usr/local/go/src/runtime/map_fast32.go:92 +0x0
  chat/server.(*Server).Run()
      /Users/sabra/go/src/chat/server/core.go:101 +0x3be

Previous read at 0x00c0001247e0 by main goroutine:
  runtime.mapiternext()
      /usr/local/go/src/runtime/map.go:851 +0x0
  chat/server.(*Server).List()
      /Users/sabra/go/src/chat/server/core.go:211 +0x167
  chat/app.(*Application).ParseInput()
      /Users/sabra/go/src/chat/app/core.go:140 +0x484
  main.main()
      /Users/sabra/go/src/chat/bin/chat/main.go:65 +0x42a

Goroutine 7 (running) created at:
  main.main()
      /Users/sabra/go/src/chat/bin/chat/main.go:46 +0x353
==================
```


### Running the Program
Starting the program can be done by entering one of two commands into terminal:
1. `./chat -port <port no>`  
2. `make run` *will automatically run on port 8888*

### Application Starting Output
```
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

Please enter a command:
```

## Programming Assignment Details
### Environment Requirements
(1) Use TCP Sockets in your peer connection implementation  
(2) Use the select() API or multi-threads  for  handling  multiple  socket  connections  
(3) Integrate **both** client-side and server-side code into **one program** and run on each peer  

### Functional Requirements
- [x] help: Display information about the available user interface options or command manual.  
- [x] myip: Display the IP address of this process.  
- [x] myport: Display the port on which this process is listening for incoming connections.  
- [x] connect {destination} {port no}:
    - [x] This  command establishes a new TCP connection to the specified {destination} at the specified {port no}. The {destination} is the IP address of the computer.  
    - [x] Any attempt to connect to an invalid IP should be rejected and suitable error message should be displayed.  
    - [x] Success or failure in connections between two peers should be indicated by both the peers using suitable messages.  
    - [x] Self-connections and duplicate connections should be flagged with suitable error messages.  
- [x] list:  
    - [x] Display a numbered list of all the connections this process is part of.  
    - [x] This numbered list will include connections initiated by this process and connections initiated by other processes.  
    - [x] The output should display the IP address and the listening port of all the peers the process is connected to.  
- [x] terminate {connection id}:  
    - [x] This command will terminate the connection listed under the specified number when 'list' is used to display all connections.  
    - [x] An error message is displayed if a valid connection does not exist.  
    - [x] If a remote machine terminates one of your connections, you should also display a message.  
- [x] send {connection id} {message}:  
    - [x] This will send the message to the host on the connection that is designated by the {connection id} when command “list” is used.  
    - [x] The message to be sent can be up-to 100 characters long, including blank spaces.  
    - [x] On successfully executing the command, the sender should display “Message sent to {connection id}” on the screen.  
    - [x] On receiving any message from the peer, the receiver should display the received message along with the sender information.  
- [x] exit:
    - [x] Close all connections and terminate this process.
    - [x] The other peers should also update their connection list by removing the peer that exits.  

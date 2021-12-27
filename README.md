# Simple Message Passing 

This is a extremelly simple example server for passing messages between applications.

This project is done only for educational purposes, for learning a bit of [Golang](https://go.dev/) and other stuffs. 
So, it is not intended to be the most powerfull or complete sofware. Its only objective is being a toy for playing with Go
and other interesting tools. Therefore, it has not got any **warranty** or support. I only make it public for others which
want to learn [Golang](https://go.dev/), so they can find some code as example. 

## How it works

The server has the option subscribe to a _topic_ and send messages to those topics.

## How to build

Import project. Run next for fixing dependencies:

```    
sergio@octubre:~/go_projects/smp-server$ go mod tidy 
```

Running from console:

- Build the server
    ```
    sergio@octubre:~/go_projects/smp-server$ go build .
    ```
- Create a user. This is an optional step. If server starts and no user exists on the system, then an arbitrary one is 
  created and printed to console.
    ```
    sergio@octubre:~/go_projects/smp-server$ ./smp --create-user sergio
    Password: ******
    ```
- Run the server
    ```
    sergio@octubre:~/proyectos/go_projects/smp-server$ ./smp
    WARN[0000] table users already exists                   
    INFO[0000] Starting the server on 0.0.0.0:1984
    ```

## How to connect

Client implementation:

- [Java client](https://github.com/scruz84/smp-java-client).

### Message types

TODO: _document the message types and formats_

Represented with one byte value:
- a: loging request
- b: loging response
- c: topic subscription
- d: send a message to a topic
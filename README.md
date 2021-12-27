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
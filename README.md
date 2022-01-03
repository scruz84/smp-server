# Simple Message Passing 

This is a extremelly simple example server for passing messages between applications.

This project is done only for educational purposes, for learning a bit of [Golang](https://go.dev/) and other stuffs. 
So, it is not intended to be the most powerfull or complete sofware. Its only objective is being a toy for playing with Go
and other interesting tools. Therefore, it has not got any warranty or support. I only make it public for others which
want to learn [Golang](https://go.dev/), so they can find some code as example. 

## How it works

The server has the option subscribe to a _topic_ and send messages to those topics.

## How to build

Import project. Run next for fixing dependencies:

``` shell
sergio@octubre:~/go_projects/smp-server$ go mod tidy 
```
## How to run

### Running from console:

- Build the server
``` shell
sergio@octubre:~/go_projects/smp-server$ go build .
```

- Create a user. This is an optional step. If server starts and no user exists on the system, then an arbitrary one is 
  created and printed to console.
``` shell
sergio@octubre:~/go_projects/smp-server$ ./smp --create-user sergio
Password: ******
```

- Run the server
``` shell
sergio@octubre:~/proyectos/go_projects/smp-server$ ./smp
WARN[0000] table users already exists                   
INFO[0000] Starting the server on 0.0.0.0:1984
```
  
### Running with docker

Repository located at [smp-server](https://hub.docker.com/r/scruz84/smp-server). Execute the container as follows. If there is not an initial users database, the server will create a default user and print the password.
``` shell
sergio@octubre:~/go_projects/smp-server$ docker run --rm -p 1984:1984 scruz84/smp-server:latest
time="2021-12-28T10:12:49Z" level=info msg="Starting the server on 0.0.0.0:1984"
time="2021-12-28T10:12:49Z" level=info msg="Initial user/password. Save them! smp-admin/4c9d269b-e"
```

For initializing the with a different user, execute like this:
``` shell
sergio@octubre:/smp-server$ docker run --rm -it -v /smp-server/data:/smp-server/data scruz84/smp-server:latest --create-user sergio
Password: ******
sergio@octubre:/smp-server$ docker run --rm -p 1984:1984 -v /smp-server/data:/smp-server/data scruz84/smp-server:latest 
time="2021-12-28T10:21:24Z" level=info msg="Starting the server on 0.0.0.0:1984"
time="2021-12-28T10:21:24Z" level=warning msg="table users already exists"
```

## How to TLS
Generate a certificate and a private key.

``` shell
sergio@octubre:~/smp-server$ mkdir config/tls -p
sergio@octubre:~/smp-server$ openssl req -new -nodes -x509 -out config/tls/server.pem -keyout config/tls/server.key -days 3650 -subj "/C=ES/ST=BIERZO/L=Ponferrada/O=SMP-ORG/OU=SMP/CN=smp-server"
```

Edit configuration file and enable TLS listener.

``` yaml
server:
  host: 0.0.0.0
  port: 1984
  tls:
    enabled: true
    tls_port: 1985
    server_key: config/tls/server.key
    server_cert: config/tls/server.pem
```

## How to connect

Client implementation:

- [Java client](https://github.com/scruz84/smp-java-client).

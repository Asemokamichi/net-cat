# NET-CAT

This project consists on recreating the NetCat in a Server-Client Architecture that can run in a server mode on a specified port listening for incoming connections, and it can be used in client mode, trying to connect to a specified port and transmitting information to the server.

## Allowed Packages
- io
- log
- os
- fmt
- net
- sync
- time
- bufio
- errors
- strings
- reflect

## Usage/Examples

### Run server locally
```bash
go run .
2023/01/30 18:20:00 Started the server at  localhost:8989
```
```bash
go run . 2525
2023/01/30 18:23:31 Started the server at  localhost:2525
```

```bash
go run . localhost 2525
[USAGE]: ./TCPChat $port
```

```bash 
go run . abc
2023/01/30 18:20:54 listen tcp: lookup tcp/asd: Servname not supported for ai_socktype
exit status 1
```
### Client example
```bash
nc $IP $port
Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    `.       | `' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     `-'       `--'
[ENTER YOUR NAME]:
```


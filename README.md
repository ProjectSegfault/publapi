# PublAPI

PublAPI is a simple API for Project Segfault's upcoming public shared system (pubnix).

## Install
An installation of Go 1.19 or higher is required.
```
git clone https://github.com/ProjectSegfault/publapi
go mod download 
go build 
./publapi
```

By default it listens to port 3000 on 127.0.0.1. You can change the port with the environment variable PUBLAPI_PORT.

## Usage
Currently, PublAPI has only two routes, /online and /signup.

| ROUTE   | TYPE | EXTRA ARGS      | DESCRIPTION                        |
|---------|------|-----------------|------------------------------------|
| /online | GET  | N/A             | Returns the number of online users |
| /signup | POST | username, email | Creates a register script and notifies admins that a new registration request was sent |


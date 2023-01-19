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

By default publapi listens on 127.0.0.1:3000. You can change the port with the environment variable PUBLAPI_PORT.

Additionally, you need to set the variable PUBLAPI_SHOUTRRRURL in order for signup notifications to work. Url Format can be found at https://containrrr.dev/shoutrrr/v0.5/services/matrix/

## Usage
Currently, PublAPI has only two routes, /users and /signup.

| ROUTE   | TYPE | EXTRA ARGS           | DESCRIPTION                     |
|---------|------|----------------------|---------------------------------|
| /users  | GET  | N/A                  | Return information about users. |
| /signup | POST | username, email, ssh | Creates a register script and notifies admins that a new registration request was sent.|

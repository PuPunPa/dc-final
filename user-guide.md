User Guide
==========
<!-- Created by Hector Jorge Morales Arch as JarlArchJernRauda and Juan Pablo as PuPumPa-->
### Requirements

- Hove installed Goland 
- Download the folder *dc-final*
- Open a terminal and type the next series of command:
$ cd <PATH>/dc-final
$ go get github.com/gorilla/mux
$ go get github.com/nanomsg/mangos
$ go get github.com/anthonynsimon/bild/tree/v0.13.0
$ go get gocv.io/x/gocv
$ go get github.com/jeasonstudio/GaussianBlur
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
$ go get -v ./...
$ export GO111MODULE=off
$ export PATH="$PATH:$(go env GOPATH)/bin"
## Using the API

# Starting the Server
For the function of this API you must have open two terminals.
The first terminal will run the server, which will receive the requests from the second terminal, therefore you must type the following command:
$ go run api.go

# Client
In the other terminal you can put 7 types of comands
   - `/ Login`
   - `/ Logout`
   - `/ Status`
   - `/ workloads`
   - `/ workloads/{workload_id}`
   - `/ images`
   - `/ images/{image_id}`
# Login

The __*Login*__ command will receive an username and a password, and returns an access token.
$ curl -u username:password http://localhost:8080/login

It displays the information as:

{
	"message": "Hi username, welcome to the DPIP System",
	"token" "OjIE89GzFw"
}

# Logout

The __*Logout*__ command will receive the generated token from the user and erases it.
$ curl -H "Authorization: Bearer <ACCESS_TOKEN>" http://localhost:8080/logout

It displays the information as:

{
	"message": "Bye username, your token has been revoked"
}

# Status

The __*Status*__ command will receive the token of the user and show a message with the time.
$ curl -H "Authorization: Bearer <ACCESS_TOKEN>" http://localhost:8080/status

It displays the information as:

{
	"message": "Hi username, the DPIP System is Up and Running"
	"time": "2015-03-07 11:06:39"
}

# Workloads
The __*Workloads*__ command will receive the token of the user and show a message with the time.
$ curl -H "Authorization: Bearer <ACCESS_TOKEN>" http://localhost:8080/status

It displays the information as:

{
	"message": "Hi username, the DPIP System is Up and Running"
	"time": "2015-03-07 11:06:39"
}

# Workloads/{workload_id}
The __*Workloads/{workload_id}*__ command will receive the token of the user and show a message with the time.
$ curl -H "Authorization: Bearer <ACCESS_TOKEN>" http://localhost:8080/status

It displays the information as:

{
	"message": "Hi username, the DPIP System is Up and Running"
	"time": "2015-03-07 11:06:39"
}

# Images
The __*Images*__ command will receive the token of the user and show a message with the time.
$ curl -H "Authorization: Bearer <ACCESS_TOKEN>" http://localhost:8080/status

It displays the information as:

{
	"message": "Hi username, the DPIP System is Up and Running"
	"time": "2015-03-07 11:06:39"
}

# images/{image_id}
The __*images/{image_id}*__ command will receive the token of the user and show a message with the time.
$ curl -H "Authorization: Bearer <ACCESS_TOKEN>" http://localhost:8080/status

It displays the information as:

{
	"message": "Hi username, the DPIP System is Up and Running"
	"time": "2015-03-07 11:06:39"
}

## Controller


## Scheduler

## Worker

$go run main.go --controller <host>:<port> --worker-name <worker_name> --tags <tag1>,<tag2>
## Documentation


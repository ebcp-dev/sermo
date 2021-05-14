## We specify the base image we need for our
## go application
FROM golang:1.16.4-alpine3.13
## We create an /app directory within our
## image that will hold our application source
## files
RUN mkdir /src
## We copy everything in the root directory
## into our /app directory
ADD . /src
## We specify that we now wish to execute 
## any further commands inside our /app
## directory
WORKDIR /src
## Add this go mod download command to pull in any dependencies
RUN go mod download
## we run go build to compile the binary
## executable of our Go program
RUN go build -o main .
## Our start command which kicks off
## our newly created binary executable
CMD ["/src/main"]
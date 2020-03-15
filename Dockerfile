FROM golang:1.14-alpine

WORKDIR /app

RUN apk add --update git
RUN go get github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon -build="go build -o .build/gonnect ./cmd/gonnect"

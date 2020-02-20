FROM golang:1.12.0-alpine3.9

WORKDIR /app

RUN apk add --update git
RUN go get github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon -build="go build -o .build/main ." -command=".build/main"

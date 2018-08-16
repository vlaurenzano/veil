# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

RUN go get github.com/vlaurenzano/veil

RUN go  get github.com/go-sql-driver/mysql

ENTRYPOINT /go/bin/veil

EXPOSE 8080
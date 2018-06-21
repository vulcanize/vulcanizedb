FROM golang:1.10.3-alpine3.7

RUN apk add --no-cache make gcc musl-dev

ADD . /go/src/github.com/vulcanize/vulcanizedb
WORKDIR /go/src/github.com/vulcanize/vulcanizedb
RUN go build -o /app main.go

ENTRYPOINT ["/app"]

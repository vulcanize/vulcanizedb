FROM golang:alpine as builder
RUN apk --update --no-cache add make git g++

# Build statically linked vDB binary (wonky path because of Dep)
RUN mkdir -p /go/src/github.com/vulcanize/vulcanizedb
ADD . /go/src/github.com/vulcanize/vulcanizedb
WORKDIR /go/src/github.com/vulcanize/vulcanizedb
RUN GCO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' .

# Build migration tool
RUN go get -u -d github.com/pressly/goose/cmd/goose
WORKDIR /go/src/github.com/pressly/goose/cmd/goose
RUN GCO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -tags='no_mysql no_sqlite' -o goose

# Second stage
FROM alpine
RUN apk --no-cache --update add openrc
RUN mkdir -p /run/openrc/ && touch /run/openrc/softlevel
COPY --from=builder /go/src/github.com/vulcanize/vulcanizedb/vulcanizedb /app/vulcanizedb
COPY --from=builder /go/src/github.com/vulcanize/vulcanizedb/environments/staging.toml /app/environments/
COPY --from=builder /go/src/github.com/vulcanize/vulcanizedb/dockerfiles/lightSync-service /etc/init.d/lightSync
COPY --from=builder /go/src/github.com/vulcanize/vulcanizedb/dockerfiles/continuousLogSync-service /etc/init.d/continuousLogSync
COPY --from=builder /go/src/github.com/vulcanize/vulcanizedb/dockerfiles/startup_script.sh /app/
COPY --from=builder /go/src/github.com/vulcanize/vulcanizedb/db/migrations/* /app/
COPY --from=builder /go/src/github.com/pressly/goose/cmd/goose/goose /app/goose

WORKDIR /app
RUN rc-status
CMD ["./goose"]

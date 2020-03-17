FROM golang:alpine

RUN apk --update --no-cache add make git g++ linux-headers
# DEBUG
RUN apk add busybox-extras

# this is probably a noob move, but I want apk from alpine for the above but need to avoid Go 1.13 below as this error still occurs https://github.com/ipfs/go-ipfs/issues/6603
FROM golang:1.12.4 as builder

# Get and build vulcanizedb
ADD . /go/src/github.com/vulcanize/vulcanizedb
WORKDIR /go/src/github.com/vulcanize/vulcanizedb
RUN GO111MODULE=on GCO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o vulcanizedb .

# Get and build vulcanize's go-ipfs fork
RUN go get -u -d github.com/ipfs/go-ipfs
WORKDIR /go/src/github.com/ipfs/go-ipfs
RUN git remote add vulcanize https://github.com/vulcanize/go-ipfs.git
RUN git fetch vulcanize
RUN git checkout -b pg_ipfs vulcanize/postgres_update
RUN GO111MODULE=on GCO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o ipfs ./cmd/ipfs

# Build migration tool
RUN go get -u -d github.com/pressly/goose/cmd/goose
WORKDIR /go/src/github.com/pressly/goose/cmd/goose
RUN GCO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -tags='no_mysql no_sqlite' -o goose .

WORKDIR /go/src/github.com/vulcanize/vulcanizedb

# app container
FROM alpine

ARG USER
ARG CONFIG_FILE
ARG EXPOSE_PORT_1
ARG EXPOSE_PORT_2

RUN adduser -Du 5000 $USER
WORKDIR /app
RUN chown $USER /app
USER $USER

# chown first so dir is writable
# note: using $USER is merged, but not in the stable release yet
COPY --chown=5000:5000 --from=builder /go/src/github.com/vulcanize/vulcanizedb/$CONFIG_FILE config.toml
COPY --chown=5000:5000 --from=builder /go/src/github.com/vulcanize/vulcanizedb/dockerfiles/super_node/startup_script.sh .


# keep binaries immutable
COPY --from=builder /go/src/github.com/vulcanize/vulcanizedb/vulcanizedb vulcanizedb
COPY --from=builder /go/src/github.com/pressly/goose/cmd/goose/goose goose
COPY --from=builder /go/src/github.com/vulcanize/vulcanizedb/db/migrations migrations/vulcanizedb
COPY --from=builder /go/src/github.com/ipfs/go-ipfs/ipfs ipfs
# XXX dir is already writeable RUN touch vulcanizedb.log

EXPOSE $EXPOSE_PORT_1
EXPOSE $EXPOSE_PORT_2

CMD ["./startup_script.sh"]

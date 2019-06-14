FROM golang:alpine as builder

RUN apk --update --no-cache add dep make git g++ linux-headers
# DEBUG
RUN apk add busybox-extras

# Get and build vulcanizedb syncAndPublish fork
RUN go get -u -d github.com/vulcanize/vulcanizedb
WORKDIR /go/src/github.com/vulcanize/vulcanizedb
RUN git checkout syncAndPublish
RUN dep ensure
RUN GCO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o vulcanizedb .

# Get and build vulcanize's go-ipfs fork
RUN go get -u -d github.com/ipfs/go-ipfs
WORKDIR /go/src/github.com/ipfs/go-ipfs
RUN git remote add vulcanize https://github.com/vulcanize/go-ipfs.git
RUN git fetch vulcanize
RUN git checkout -b pg_ipfs vulcanize/postgres_update
RUN dep ensure
RUN GCO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o ipfs ./cmd/ipfs

# Get and build vulcanize's geth fork
RUN go get -u -d github.com/ethereum/go-ethereum
WORKDIR /go/src/github.com/ethereum/go-ethereum
RUN git remote add vulcanize https://github.com/vulcanize/go-ethereum.git
RUN git fetch vulcanize
RUN git checkout -b statediff_geth vulcanize/rpc_statediffing
RUN GCO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o geth ./cmd/geth

# Build migration tool
RUN go get -u -d github.com/pressly/goose/cmd/goose
WORKDIR /go/src/github.com/pressly/goose/cmd/goose
RUN GCO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -tags='no_mysql no_sqlite' -o goose .

WORKDIR /go/src/github.com/vulcanize/vulcanizedb

# app container
FROM alpine
WORKDIR /app

ARG USER
ARG config_file=environments/syncPublishScreenAndServe.toml
ARG vdb_dbname="vulcanize_public"
ARG vdb_hostname="localhost"
ARG vdb_port="5432"
ARG vdb_user="postgres"
ARG vdb_password

# setup environment
ENV VDB_PG_NAME="$vdb_dbname"
ENV VDB_PG_HOSTNAME="$vdb_hostname"
ENV VDB_PG_PORT="$vdb_port"
ENV VDB_PG_USER="$vdb_user"
ENV VDB_PG_PASSWORD="$vdb_password"

RUN adduser -D 5000 $USER
USER $USER

# chown first so dir is writable
# note: using $USER is merged, but not in the stable release yet
COPY --chown=5000:5000 --from=builder /go/src/github.com/vulcanize/vulcanizedb/$config_file config.toml
COPY --chown=5000:5000 --from=builder /go/src/github.com/vulcanize/vulcanizedb/dockerfiles/seed_node/startup_script.sh .

# keep binaries immutable
COPY --from=builder /go/src/github.com/vulcanize/vulcanizedb/vulcanizedb vulcanizedb
COPY --from=builder /go/src/github.com/pressly/goose/cmd/goose/goose goose
COPY --from=builder /go/src/github.com/vulcanize/vulcanizedb/db/migrations migrations/vulcanizedb
COPY --from=builder /go/src/github.com/ipfs/go-ipfs/ipfs ipfs
COPY --from=builder /go/src/github.com/ethereum/go-ethereum/geth geth

# DEBUG
COPY --from=builder /usr/bin/telnet /bin/telnet

EXPOSE 80

CMD ["./startup_script.sh"]

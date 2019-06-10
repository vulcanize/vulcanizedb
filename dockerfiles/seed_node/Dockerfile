FROM golang:alpine as builder

RUN apk --update --no-cache add make git g++
# DEBUG
RUN apk add busybox-extras

# Build statically linked vDB binary (wonky path because of Dep)
WORKDIR /go/src/github.com/vulcanize/vulcanizedb
ADD . .
RUN GCO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' .

# Build statically linked IPFS binary (statically linked and wonky path because we need to use custom fork)
WORKDIR /go/src/github.com/ipfs/go-ipfs
ADD . .
RUN GCO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' .

# Build statically linked Geth binary (statically linked and wonky path because we need to use custom fork)
WORKDIR /go/src/github.com/ethereum/go-ethereum
ADD . .
RUN GCO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' .

# Build migration tool
RUN go get -u -d github.com/pressly/goose/cmd/goose
WORKDIR /go/src/github.com/pressly/goose/cmd/goose
RUN GCO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -tags='no_mysql no_sqlite' -o goose

WORKDIR /go/src/github.com/vulcanize/vulcanizedb

# app container
FROM alpine
WORKDIR /app

ARG USER
ARG config_file=environments/example.toml
ARG vdb_command=syncPublishScreenAndServe
ARG vdb_pg_connect="postgres://$USER@/vulcanize_public?sslmode=disable"
ARG vdb_dbname="vulcanize_public"

# setup environment
ENV VDB_COMMAND="$vdb_command"
ENV VDB_PG_CONNECT="$vdb_pg_connect"

RUN adduser -Du 5000 $USER
USER $USER

# chown first so dir is writable
# note: using $USER is merged, but not in the stable release yet
COPY --chown=5000:5000 --from=builder /go/src/github.com/vulcanize/vulcanizedb/$config_file config.toml
COPY --chown=5000:5000 --from=builder /go/src/github.com/vulcanize/vulcanizedb/dockerfiles/seed_node/startup_script.sh .

# keep binaries immutable
COPY --from=builder /go/src/github.com/vulcanize/vulcanizedb/vulcanizedb vulcanizedb
COPY --from=builder /go/src/github.com/pressly/goose/cmd/goose/goose goose
COPY --from=builder /go/src/github.com/vulcanize/vulcanizedb/db/migrations migrations/vulcanizedb
COPY --from=builder /go/src/github.com/ipfs/go-ipfs ipfs
COPY --from=builder /go/src/github.com/ethereum/go-ethereum geth

# DEBUG
COPY --from=builder /usr/bin/telnet /bin/telnet

CMD ["./startup_script.sh"]

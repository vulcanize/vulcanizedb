FROM golang:alpine as builder

RUN apk --update --no-cache add make git g++ linux-headers
# DEBUG
RUN apk add busybox-extras

# build migration tool
RUN go get -u -d github.com/pressly/goose/cmd/goose
WORKDIR /go/src/github.com/pressly/goose/cmd/goose
RUN GCO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -tags='no_mysql no_sqlite' -o goose

# build statically linked vDB binary (wonky path because of Dep)
WORKDIR /go/src/github.com/vulcanize/vulcanizedb
ADD . .
RUN GCO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' .


# app container
FROM alpine
WORKDIR /app

RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true

ARG USER
ARG config_file=environments/example.toml
ARG vdb_command=headerSync
ARG vdb_pg_host="localhost"
ARG vdb_pg_port="5432"
ARG vdb_dbname="vulcanize_public"
ARG vdb_pg_connect="postgres://$USER@$vdb_pg_host:$vdb_pg_port/$vdb_dbname?sslmode=disable"

# setup environment
ENV VDB_COMMAND="$vdb_command"
ENV VDB_PG_CONNECT="$vdb_pg_connect"

RUN adduser -Du 5000 $USER
USER $USER

# chown first so dir is writable
# note: using $USER is merged, but not in the stable release yet
COPY --chown=5000:5000 --from=builder /go/src/github.com/vulcanize/vulcanizedb/$config_file config.toml
COPY --chown=5000:5000 --from=builder /go/src/github.com/vulcanize/vulcanizedb/dockerfiles/startup_script.sh .

# keep binaries immutable
COPY --from=builder /go/src/github.com/vulcanize/vulcanizedb/vulcanizedb vulcanizedb
COPY --from=builder /go/src/github.com/pressly/goose/cmd/goose/goose goose
COPY --from=builder /go/src/github.com/vulcanize/vulcanizedb/plugins/ .
COPY --from=builder /go/src/github.com/vulcanize/vulcanizedb/db/migrations migrations/vulcanizedb

# DEBUG
COPY --from=builder /usr/bin/telnet /bin/telnet

CMD ["./startup_script.sh"]

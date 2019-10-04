FROM golang:alpine as builder

RUN apk --update --no-cache add make git g++ linux-headers
# DEBUG
RUN apk add busybox-extras

# Build statically linked vDB binary (wonky path because of Dep)
WORKDIR /go/src/github.com/vulcanize/vulcanizedb
ADD . .
RUN GO111MODULE=on GCO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' .

# Build migration tool
WORKDIR /go
RUN go get -u -d github.com/pressly/goose/cmd/goose
WORKDIR /go/src/github.com/pressly/goose/cmd/goose
RUN GCO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -tags='no_mysql no_sqlite' -o goose

WORKDIR /go/src/github.com/vulcanize/vulcanizedb


# app container
FROM alpine
WORKDIR /app

RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true

ARG USER
ARG config_file=environments/moloch.toml

RUN adduser -Du 5000 $USER
USER $USER

# chown first so dir is writable
# note: using $USER is merged, but not in the stable release yet
COPY --chown=5000:5000 --from=builder /go/src/github.com/vulcanize/vulcanizedb/$config_file config.toml
COPY --chown=5000:5000 --from=builder /go/src/github.com/vulcanize/vulcanizedb/dockerfiles/startup_script.sh .

# keep binaries immutable
COPY --from=builder /go/src/github.com/vulcanize/vulcanizedb/vulcanizedb vulcanizedb
COPY --from=builder /go/src/github.com/pressly/goose/cmd/goose/goose goose
COPY --from=builder /go/src/github.com/vulcanize/vulcanizedb/db/migrations db/migrations

# DEBUG
COPY --from=builder /usr/bin/telnet /bin/telnet

CMD ["./startup_script.sh"]

FROM golang:alpine as builder

RUN apk --update --no-cache add g++ git linux-headers

# Build migration tool
WORKDIR /go
RUN GO111MODULE=auto go get -u -d github.com/pressly/goose/cmd/goose
WORKDIR /go/src/github.com/pressly/goose/cmd/goose
RUN GO111MODULE=auto go build -a -ldflags '-s' -tags='no_mysql no_sqlite' -o goose

ENV GO111MODULE on

WORKDIR /vulcanizedb
COPY . .
RUN go build


FROM golang:alpine

WORKDIR /app

# add certificates for node requests via https
RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true

# setup environment
ENV GO111MODULE on

# Direct logs to stdout for docker log driver
RUN ln -sf /dev/stdout /app/vulcanizedb.log

# add required files
COPY --from=builder /vulcanizedb/vulcanizedb .
COPY --from=builder /vulcanizedb/dockerfiles/header_sync/startup_script.sh .
COPY --from=builder /vulcanizedb/db/migrations/* db/migrations/
COPY --from=builder /go/src/github.com/pressly/goose/cmd/goose/goose goose
COPY --from=builder /vulcanizedb/dockerfiles/wait-for-it.sh .

# need to execute with a shell to access env variables
CMD ["./startup_script.sh"]

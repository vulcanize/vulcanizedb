FROM golang:alpine as builder

RUN apk --update --no-cache add g++ git linux-headers

ENV GO111MODULE on

WORKDIR /vulcanizedb
COPY . .
RUN go build
# I think this should build the same version of goose that was already downloaded as part of go.mod,
# but we should validate that.
RUN go install -tags='no_mysql no_sqlite3 no_mssql no_redshift' github.com/pressly/goose/cmd/goose
RUN find / -name goose -print
RUN which goose

FROM golang:alpine

WORKDIR /app

# add certificates for node requests via https
# bash for wait-for-it.sh
RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
	bash \
        && update-ca-certificates 2>/dev/null || true

# setup environment
ENV GO111MODULE on

# Direct logs to stdout for docker log driver
RUN ln -sf /dev/stdout /app/vulcanizedb.log

# add required files
COPY --from=builder /vulcanizedb/vulcanizedb .
COPY --from=builder /vulcanizedb/dockerfiles/header_sync/startup_script.sh .
COPY --from=builder /vulcanizedb/db/migrations/* db/migrations/
COPY --from=builder /go/bin/goose goose
# needed for waiting until postgres is ready before starting from docker-compose
COPY --from=builder /vulcanizedb/dockerfiles/wait-for-it.sh .

HEALTHCHECK CMD test -f /tmp/header_sync_health_check

# need to execute with a shell to access env variables
CMD ["./startup_script.sh"]

FROM golang:alpine

RUN apk --update --no-cache add make git g++ linux-headers
# DEBUG
RUN apk add busybox-extras

# this is probably a noob move, but I want apk from alpine for the above but need to avoid Go 1.13 below as this error still occurs https://github.com/ipfs/go-ipfs/issues/6603
FROM golang:1.12.4 as builder

# Get and build vulcanizedb
ADD . /go/src/github.com/vulcanize/vulcanizedb

# Build migration tool
RUN go get -u -d github.com/pressly/goose/cmd/goose
WORKDIR /go/src/github.com/pressly/goose/cmd/goose
RUN GCO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -tags='no_mysql no_sqlite' -o goose .

WORKDIR /go/src/github.com/vulcanize/vulcanizedb

# app container
FROM alpine

ARG USER

RUN adduser -Du 5000 $USER
WORKDIR /app
RUN chown $USER /app
USER $USER

# chown first so dir is writable
# note: using $USER is merged, but not in the stable release yet
COPY --chown=5000:5000 --from=builder /go/src/github.com/vulcanize/vulcanizedb/dockerfiles/migrations/startup_script.sh .


# keep binaries immutable
COPY --from=builder /go/src/github.com/pressly/goose/cmd/goose/goose goose
COPY --from=builder /go/src/github.com/vulcanize/vulcanizedb/db/migrations migrations/vulcanizedb
# XXX dir is already writeable RUN touch vulcanizedb.log

CMD ["./startup_script.sh"]
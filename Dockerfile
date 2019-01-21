FROM golang:alpine as builder
RUN apk --update --no-cache add make git g++

# Build statically linked vDB binary (wonky path because of Dep)
RUN mkdir -p /go/src/github.com/vulcanize/vulcanizedb
ADD . /go/src/github.com/vulcanize/vulcanizedb
WORKDIR /go/src/github.com/vulcanize/vulcanizedb
RUN GCO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' .

# Second stage
FROM scratch
COPY --from=builder /go/src/github.com/vulcanize/vulcanizedb/vulcanizedb /app/vulcanizedb
WORKDIR /app
CMD ["./vulcanizedb", "--help"]

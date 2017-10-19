# Vulcanize DB

[![Build Status](https://travis-ci.com/8thlight/vulcanizedb.svg?token=GKv2Y33qsFnfYgejjvYx&branch=master)](https://travis-ci.com/8thlight/vulcanizedb)

## Development Setup

By default, `go get` does not work for private GitHub repos. This will fix that.
1. `git config --global url."git@github.com:".insteadOf "https://github.com/"`
2. `go get github.com/8thlight/vulcanizedb`

## Running the Tests

`go test ./...`

# Vulcanize DB

[![Build Status](https://travis-ci.com/8thlight/vulcanizedb.svg?token=GKv2Y33qsFnfYgejjvYx&branch=master)](https://travis-ci.com/8thlight/vulcanizedb)

## Development Setup

By default, `go get` does not work for private GitHub repos. This will fix that.
1. `git config --global url."git@github.com:".insteadOf "https://github.com/"`
2. `go get github.com/8thlight/vulcanizedb`

## Running the Tests

### Integration Test

In order to run the integration tests, you will need to run them against a real blockchain. Here are steps to create a local, private blockchain.

1. Run `./scripts/setup` to create a private blockchain with a new account.
    * This will result in a warning.
2. Run `./scripts/start_private_blockchain` as a separate process.
3. `go test ./...`

### Unit Tests

`go test ./core`
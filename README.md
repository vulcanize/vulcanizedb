# Vulcanize DB

[![Build Status](https://travis-ci.com/8thlight/vulcanizedb.svg?token=GKv2Y33qsFnfYgejjvYx&branch=master)](https://travis-ci.com/8thlight/vulcanizedb)

## Development Setup

### Dependencies

 - Go 1.9+
 - https://github.com/golang/dep
    - `go get -u github.com/golang/dep/cmd/dep`
 - Postgres 10

### Cloning the Repository

1. `git config --global url."git@github.com:".insteadOf "https://github.com/"`
    - By default, `go get` does not work for private GitHub repos. This will fix that.
2. `go get github.com/8thlight/vulcanizedb`
3. `go get github.com/ethereum/go-ethereum`
    - This will take a while and gives poor indication of progress.
4. `go install github.com/ethereum/go-ethereum/cmd/geth`
5. `cd $GOPATH/src/github.com/8thlight/vulcanizedb`
6. `dep ensure`

### Setting up the Development Database

1. Install Postgres
2. Create a superuser for yourself and make sure `psql --list` works without prompting for a password.
3. `go get -u -d github.com/mattes/migrate/cli github.com/lib/pq`
4. `go build -tags 'postgres' -o /usr/local/bin/migrate github.com/mattes/migrate/cli`
5. `createdb vulcanize`
6. `cd $GOPATH/src/github.com/8thlight/vulcanizedb`
7.  `./scripts/migrate`

Adding a new migration: `./scripts/create_migration <migration-name>`

### Creating/Using a Private Blockchain

Syncing the public blockchain takes many hours for the initial sync and will download 20+ GB of data.
Here are some instructions for creating a private blockchain that does not depend on having a network connection.

1. Run `./scripts/setup` to create a private blockchain with a new account.
    * This will result in a warning.
2. Run `./scripts/start_private_blockchain`.

### Connecting to the Public Blockchain

`./scripts/start_blockchain`

### IPC File Paths

The default location for Ethereum is:
 - `$HOME/Library/Ethereum` for Mac
 - `$HOME/.ethereum` for Ubuntu
 - `$GOPATH/src/gihub.com/8thlight/vulcanizedb/test_data_dir/geth.ipc` for private blockchain.

**Note the location of the ipc file is outputted when you connect to a blockchain. It is needed to start the listener below**

## Running Listener

1. Start a blockchain.
2. In a separate terminal start listener (ipcDir location)
    - `go run main.go --ipcPath /path/to/file.ipc`

## Running the Tests

### Integration Test

In order to run the integration tests, you will need to run them against a real blockchain.

1. Run `./scripts/start_private_blockchain` as a separate process.
2. `go test ./...`

### Unit Tests

1. `go test ./core`

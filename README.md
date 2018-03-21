# Vulcanize DB

[![Build Status](https://travis-ci.org/vulcanize/VulcanizeDB.svg?branch=master)](https://travis-ci.org/vulcanize/VulcanizeDB)

### Dependencies

 - Go 1.9+
 - Postgres 10 
 - Ethereum Node
   - [Go Ethereum](https://ethereum.github.io/go-ethereum/downloads/) (1.8+)
   - [Parity 1.8.11+](https://github.com/paritytech/parity/releases)
   
### Installation 
`go get github.com/vulcanize/vulcanizedb`

### Setting up the Databases

1. Install Postgres
1. Create a superuser for yourself and make sure `psql --list` works without prompting for a password.
1. `createdb vulcanize_private`
1. `cd $GOPATH/src/github.com/vulcanize/vulcanizedb`
1. Import the schema: `psql vulcanize_private < db/schema.sql`

   or run the migrations: `make migrate HOST_NAME=localhost NAME=vulcanize_public PORT=5432`
    * See below for configuring additional environments
    
### IPC File Paths

The default location for Ethereum is:
 - `$HOME/Library/Ethereum` for Mac
 - `$HOME/.ethereum` for Ubuntu
 - `$GOPATH/src/gihub.com/vulcanize/vulcanizedb/test_data_dir/geth.ipc` for private node.

**Note the location of the ipc file is printed to the console when you start geth. It is needed to for configuration**

## Start syncing with postgres
1. Start geth node (**if fast syncing wait for geth to finsh initial sync**)
1. In a separate terminal start vulcanize_db
    - `vulcanizedb sync --config <config.toml> --starting-block-number <block-number>`
    
   * see `./environments` for example config 

## Running the Tests

### Unit Tests

1. `go test ./pkg/...`

### Integration Test

1. Setup a test database and import the schema: 

   `createdb vulcanize_private`
   
    `psql vulcanize_private < db/schema.sql`
1. `go test ./...` to run all tests.

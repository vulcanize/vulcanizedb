# Vulcanize DB

[![Build Status](https://travis-ci.com/8thlight/vulcanizedb.svg?token=3psFYN2533rYjhRbvjte&branch=master)](https://travis-ci.com/8thlight/vulcanizedb)

### Dependencies

 - Go 1.9+
 - Postgres 10 
 - Go Ethereum 
    - https://ethereum.github.io/go-ethereum/downloads/ 
 
### Installation 
    go get github.com/vulcanize/vulcanizedb
    
### Setting up the Databases

1. Install Postgres
2. Create a superuser for yourself and make sure `psql --list` works without prompting for a password.
3. `createdb vulcanize_private`
4. `cd $GOPATH/src/github.com/vulcanize/vulcanizedb`
5. Import the schema: `psql vulcanize_private < db/schema.sql`

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
2. In a separate terminal start vulcanize_db
    - `vulcanizedb sync --config <config.toml> --starting-block-number <block-number>`
    
   * see `./environments` for example config 

## Watch specific events
1. Start geth 
2. In a separate terminal start vulcanize_db
    - `vulcanizedb sync --config <config.toml> --starting-block-number <block-number>`
3. Create event filter 
    - `vulcanizedb addFilter --config <config.toml> --filter-filepath <filter.json>`
   * see `./filters` for example filter 
4. The filters are tracked in the `log_filters` table and the filtered events 
will show up in the `watched_log_events` view
     
## Running the Tests

### Unit Tests

1. `go test ./pkg/...`

### Integration Test

In order to run the integration tests, you will need to run them against a real node. At the moment the integration tests require [Geth v1.7.2](https://ethereum.github.io/go-ethereum/downloads/) as they depend on the `--dev` mode, which changed in v1.7.3 

1. Run `make startprivate` in a separate terminal
2. Setup a test database and import the schema: 

   `createdb vulcanize_private`
   
    `psql vulcanize_private < db/schema.sql`
3. `go test ./...` to run all tests.

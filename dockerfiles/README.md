S
`Dockerfile` will build an alpine image containing:
- vDB as a binary with runtime deps statically linked: `/app/vulcanizedb`
- The migration tool goose: `/app/goose`
- Two services for running `headerSync` and `continuousLogSync`, started with the default configuration `environments/staging.toml`.

By default, vDB is configured towards the Kovan deploy. The configuration values can be overridden using environment variables, using the same hierarchical naming pattern but in CAPS and using underscores. For example, the contract address for the `Pit` can be set with the variable `CONTRACT_ADDRESS_PIT="0x123..."`.


## To build the container:

By default will use `environments/example.toml` which will try to connect to a local postgres database and kovan instance on port 8545.

Passing the current USER as a build arg is required, with `--build-arg USER`
```sh
export VDB_VERSION=v0.0.2
docker build -t vdb:$VDB_VERSION . --build-arg USER
```

### Build args
- `config_file` - specify the path to an alternate config file to use
- `vdb_command` - space delimited list of vulcanizedb commands to run
- `vdb_pg_connect` - goose URL for migrations

Defaults are located in the dockerfile.


## To run the docker image

```sh
export VDB_VERSION=v0.0.2
docker run --network host vdb:$VDB_VERSION
```

## To use the container:
1. Setup a postgres database with superuser `vulcanize`
2. Set the env variables `VDB_PG_NAME`, `VDB_PG_HOSTNAME`,
  `VDB_PG_PORT`, `VDB_PG_USER` & `VDB_PG_PASSWORD`
3. Run the DB migrations:
  * `./goose postgres "postgresql://$(VDB_PG_USER):$(VDB_PG_PASSWORD)@$(VDB_PG_HOSTNAME):$(VDB_PG_PORT)/$(VDB_PG_NAME)?sslmode=disable"
e`
4. Set `CLIENT_IPCPATH` to a node endpoint
5. Set the contract variables:
  * `CONTRACT_ADDRESS_[CONTRACT NAME]=0x123...`
  * `CONTRACT_ABI_[CONTRACT NAME]="ABI STRING"`
  * `CONTRACT_DEPLOYMENT-BLOCK_[CONTRACT NAME]=0` (doesn't really matter on a short chain, just avoids long unnecessary searching)
6. Start the `headerSync` and `continuousLogSync` services:
  * `./vulcanizedb headerSync --config environments/staging.toml`
  * `./vulcanizedb continuousLogSync --config environments/staging.toml`

### Automated
The steps above have been rolled into a script: `/app/startup_script.sh`, which just assumes the DB env variables have been set, and defaults the rest to Kovan according to `environments/staging.toml`. This can be called with something like:

`docker run -d -e VDB_PG_NAME=vulcanize_public -e VDB_PG_HOSTNAME=localhost -e VDB_PG_PORT=5432 -e VDB_PG_USER=vulcanize -e VDB_PG_PASSWORD=vulcanize m0ar/images:vDB`

### Logging
When running, vDB services log to `/app/vulcanizedb.log`.


S
`Dockerfile` will build an alpine image containing:
- vDB as a binary with runtime deps statically linked: `/app/vulcanizedb`
- The migration tool goose: `/app/goose`
- Two services for running `lightSync` and `continuousLogSync`, started with the default configuration `environments/staging.toml`.

By default, vDB is configured towards the Kovan deploy. The configuration values can be overridden using environment variables, using the same hierarchical naming pattern but in CAPS and using underscores. For example, the contract address for the `Pit` can be set with the variable `CONTRACT_ADDRESS_PIT="0x123..."`.

## To use the container:
1. Setup a postgres database with owner `vulcanize`
  * `vulcanize` does not _have_ to be owner, but otherwise permissions for
    connection, tables, and sequences need to be added for the `public` and
    `maker` schemas manually when the migrations have been run.
2. Set the env variables `DATABASE_NAME`, `DATABASE_HOSTNAME`,
  `DATABASE_PORT`, `DATABASE_USER` & `DATABASE_PASSWORD`
3. Run the DB migrations:
  * `./goose postgres "postgresql://$(DATABASE_USER):$(DATABASE_PASSWORD)@$(DATABASE_HOSTNAME):$(DATABASE_PORT)/$(DATABASE_NAME)?sslmode=disable"
e`
4. Set `CLIENT_IPCPATH` to a node endpoint
5. Set the contract variables:
  * `CONTRACT_ADDRESS_[CONTRACT NAME]=0x123...`
  * `CONTRACT_ABI_[CONTRACT NAME]="ABI STRING"`
  * `CONTRACT_DEPLOYMENT-BLOCK_[CONTRACT NAME]=0` (doesn't really matter on a short chain, just avoids long unnecessary searching)
6. Start the `lightSync` and `continuousLogSync` services:
  * `rc-service lightSync start`
  * `rc-service continuousLogSync start`

### Automated
The steps above have been rolled into a script: `/app/startup_script.sh`, which just assumes the DB env variables have been set, and defaults the rest to Kovan according to `environments/staging.toml`.

### Logging
When running, vDB services log to `/vulcanizedb.log`.


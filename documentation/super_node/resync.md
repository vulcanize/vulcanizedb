## VulcanizeDB Super Node Resync
The `resync` command is made available for directing the resyncing of super node data within specified ranges.
It also contains a utility for cleaning out old data, and resetting the validation level of data.

### Rational

Manual resyncing of data is useful when we want to re-validate data within specific ranges using a new source.

Cleaning out data is useful when we need to remove bad/deprecated data or prepare for breaking changes to the db schemas.

Resetting the validation level of data is useful for designating ranges of data for resyncing by an ongoing super node
backfill process.

### Command

Usage: `./vulcanizedb resync --config={config.toml}`

Configuration can also be done through CLI options and/or environmental variables.
CLI options can be found using `./vulcanizedb resync --help`.

### Config

Below is the set of universal config parameters for the resync command, in .toml form, with the respective environmental variables commented to the side.
This set of parameters needs to be set no matter the chain type.

```toml
[database]
    name     = "vulcanize_public" # $DATABASE_NAME
    hostname = "localhost" # $DATABASE_HOSTNAME
    port     = 5432 # $DATABASE_PORT
    user     = "vdbm" # $DATABASE_USER
    password = "" # $DATABASE_PASSWORD

[ipfs]
    path = "~/.ipfs" # $IPFS_PATH
    
[resync]
    chain = "ethereum" # $RESYNC_CHAIN
    type = "state" # $RESYNC_TYPE
    start = 0 # $RESYNC_START
    stop = 1000 # $RESYNC_STOP
    batchSize = 10 # $RESYNC_BATCH_SIZE
    batchNumber = 100 # $RESYNC_BATCH_NUMBER
    clearOldCache = true # $RESYNC_CLEAR_OLD_CACHE
    resetValidation = true # $RESYNC_RESET_VALIDATION
```

Additional parameters need to be set depending on the specific chain.

For Bitcoin: 

```toml
[bitcoin]
    httpPath = "127.0.0.1:8332" # $BTC_HTTP_PATH
    pass = "password" # $BTC_NODE_PASSWORD
    user = "username" # $BTC_NODE_USER
    nodeID = "ocd0" # $BTC_NODE_ID
    clientName = "Omnicore" # $BTC_CLIENT_NAME
    genesisBlock = "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f" # $BTC_GENESIS_BLOCK
    networkID = "0xD9B4BEF9" # $BTC_NETWORK_ID
```

For Ethereum:

```toml
[ethereum]
    httpPath = "127.0.0.1:8545" # $ETH_HTTP_PATH
```

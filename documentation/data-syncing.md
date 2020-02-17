# Syncing commands
These commands are used to sync raw Ethereum data into Postgres, with varying levels of data granularity.

## headerSync
Syncs block headers from a running Ethereum node into the VulcanizeDB table `headers`.
- Queries the Ethereum node using RPC calls.
- Validates headers from the last 15 blocks to ensure that data is up to date.
- Useful when you want a minimal baseline from which to track targeted data on the blockchain (e.g. individual smart
contract storage values or event logs).
- Handles chain reorgs by [validating the most recent blocks' hashes](../pkg/history/header_validator.go). If the hash is
different from what we have already stored in the database, the header record will be updated.

#### Usage
- Run: `./vulcanizedb headerSync --config <config.toml> --starting-block-number <block-number>`
- The config file must be formatted as follows, and should contain an ipc path to a running Ethereum node:
```toml
[database]
    name     = "vulcanize_public"
    hostname = "localhost"
    user     = "vulcanize"
    password = "vulcanize"
    port     = 5432

[client]
    ipcPath  = <path to a running Ethereum node>
```
- Alternatively, the ipc path can be passed as a flag instead `--client-ipcPath`.

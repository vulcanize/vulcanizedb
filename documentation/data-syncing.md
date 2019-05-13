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

## fullSync
Syncs blocks, transactions, receipts and logs from a running Ethereum node into VulcanizeDB tables named
`blocks`, `uncles`, `full_sync_transactions`, `full_sync_receipts` and `logs`. 
- Queries the Ethereum node using RPC calls.
- Validates headers from the last 15 blocks to ensure that data is up to date.
- Useful when you want to maintain a broad cache of what's happening on the blockchain.
- Handles chain reorgs by [validating the most recent blocks' hashes](../pkg/history/header_validator.go). If the hash is
different from what we have already stored in the database, the header record will be updated.

#### Usage
- Run `./vulcanizedb fullSync --config <config.toml> --starting-block-number <block-number>`
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

*Please note, that if you are fast syncing your Ethereum node, wait for the initial sync to finish.*

## coldImport
Syncs VulcanizeDB from Geth's underlying LevelDB datastore and persists Ethereum blocks, 
transactions, receipts and logs into VulcanizeDB tables named `blocks`, `uncles`, 
`full_sync_transactions`, `full_sync_receipts` and `logs` respectively.

#### Usage
1. Ensure the Ethereum node you're point at is not running, and that it has synced to the desired block height.
1. Run `./vulcanizedb coldImport --config <config.toml>`
1. Optional flags:
    - `--starting-block-number <block number>`/`-s <block number>`: block number to start syncing from
    - `--ending-block-number <block number>`/`-e <block number>`: block number to sync to
    - `--all`/`-a`: sync all missing blocks

The config file can be formatted as follows, and must contain the LevelDB path.

```toml
[database]
    name     = "vulcanize_public"
    hostname = "localhost"
    user     = "vulcanize"
    password = "vulcanize"
    port     = 5432

[client]
    leveldbpath = "/Users/user/Library/Ethereum/geth/chaindata"
```

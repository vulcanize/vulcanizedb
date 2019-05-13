# Syncing commands
These commands are used to sync raw Ethereum data into Postgres, with varying levels of data granularity.

## headerSync
Syncs block headers from a running Ethereum node into the VulcanizeDB table `headers`.
- Queries the Ethereum node using RPC calls.
- Validates headers from the last 15 blocks to ensure that data is up to date.
- Useful when you want a minimal baseline from which to track targeted data on the blockchain (e.g. individual smart contract storage values or event logs).

##### Usage
1. Start Ethereum node.
1. In a separate terminal start VulcanizeDB:
`./vulcanizedb headerSync --config <config.toml> --starting-block-number <block-number>`

## sync
Syncs blocks, transactions, receipts and logs from a running Ethereum node into VulcanizeDB tables named
`blocks`, `uncles`, `full_sync_transactions`, `full_sync_receipts` and `logs`. 
- Queries the Ethereum node using RPC calls.
- Validates headers from the last 15 blocks to ensure that data is up to date.
- Useful when you want to maintain a broad cache of what's happening on the blockchain.

##### Usage
1. Start Ethereum node (**if fast syncing your Ethereum node, wait for initial sync to finish**).
1. In a separate terminal start VulcanizeDB:
`./vulcanizedb sync --config <config.toml> --starting-block-number <block-number>`

## coldImport
Syncs VulcanizeDB from Geth's underlying LevelDB datastore and persists Ethereum blocks, 
transactions, receipts and logs into VulcanizeDB tables named `blocks`, `uncles`, 
`full_sync_transactions`, `full_sync_receipts` and `logs` respectively.

##### Usage
1. Assure node is not running, and that it has synced to the desired block height.
1. Start vulcanize_db
   - `./vulcanizedb coldImport --config <config.toml>`
1. Optional flags:
    - `--starting-block-number <block number>`/`-s <block number>`: block number to start syncing from
    - `--ending-block-number <block number>`/`-e <block number>`: block number to sync to
    - `--all`/`-a`: sync all missing blocks

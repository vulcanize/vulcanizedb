# Syncing commands
These commands are used to sync raw Ethereum data into Postgres.

## lightSync
Syncs VulcanizeDB with the configured Ethereum node, populating only block headers.
This command is useful when you want a minimal baseline from which to track targeted data on the blockchain (e.g. individual smart contract storage values or event logs).
1. Start Ethereum node
1. In a separate terminal start VulcanizeDB:
    - `./vulcanizedb lightSync --config <config.toml> --starting-block-number <block-number>`
    
## sync
Syncs VulcanizeDB with the configured Ethereum node, populating blocks, transactions, receipts, and logs.
This command is useful when you want to maintain a broad cache of what's happening on the blockchain.
1. Start Ethereum node (**if fast syncing your Ethereum node, wait for initial sync to finish**)
1. In a separate terminal start VulcanizeDB:
    - `./vulcanizedb sync --config <config.toml> --starting-block-number <block-number>`

## coldImport
Sync VulcanizeDB from the LevelDB underlying a Geth node.
1. Assure node is not running, and that it has synced to the desired block height.
1. Start vulcanize_db
   - `./vulcanizedb coldImport --config <config.toml>`
1. Optional flags:
    - `--starting-block-number <block number>`/`-s <block number>`: block number to start syncing from
    - `--ending-block-number <block number>`/`-e <block number>`: block number to sync to
    - `--all`/`-a`: sync all missing blocks

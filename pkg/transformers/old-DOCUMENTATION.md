The main goal of creating a transformer is to fetch specific log events from Ethereum, convert/decode them into usable data and then persist them to VulcanizeDB. For Maker there are two main types of log events that we're tracking: custom events that are defined in the contract solidity code, and LogNote events which utilize the [DSNote library](https://github.com/dapphub/ds-note). The transformer process for each of these different log types is the same, except for the converting process, as denoted below.

## Creating a Transformer for custom events (i.e. FlopKick)
To illustrate how to create a custom log event transformer we'll use the Kick event defined in [flop.sol](https://github.com/makerdao/dss/blob/master/src/flop.sol) as an example.

1. Get an example FlopKick log event either from mainnet (if the
   contract has already been deployed), from the Kovan testnet, or by
deploying the contract to a local chain and emitting the event manually.
We will use the example log event to test drive converting the log to a
FlopKick database model.

1. Fetch the appropriate logs from the chain.
  - Most transformers use `shared.LogFetcher`
  - Price Feeds

1. Convert the raw log into a database model.
  - For Custom Events, such as FlopKick
    1. Create a converter to convert the raw log into a Go structure.
      - We've been using [go-ethereum's abigen tool](https://github.com/ethereum/go-ethereum/tree/master/cmd/abigen) to get the contract's ABI, and a Go struct that represents the event log. We will unpack the raw logs into this struct.
        - To use abigen: `abigen --sol flip.sol --pkg flip --out {/path/to/output_file}`
          - sol: this is the path to the solidity contract
          - pkg: a package name for the generated Go code
          - out: the file path for the generated Go code (optional)
        - the output for `flop.sol` will include the FlopperAbi and the FlopperKick struct:
        ```go
            type FlopperKick struct {
              Id  *big.Int
              Lot *big.Int
              Bid *big.Int
              Gal common.Address
              End *big.Int
              Raw types.Log
            }
        ```
      - Using go-ethereum's `contract.UnpackLog` method we can unpack the raw log into the FlopperKick struct (which we're referring to as the `entity`).
        - See the `ToEntity` method in `pkg/transformers/flop_kick/converter`.
        - The unpack method will not add the `Raw` or `TransactionIndex` values to the entity struct - both of these values are accessible from the entity.
    1. Then convert the entity into a database model. See the `ToModel` method in `pkg/transformers/flop_kick/converter`.

  - For LogNote Events, such as tend.
    - Since LogNote events are a generic structure, they depend on the
      method signature of the method that is calling them. For example,
when the `tend` method is called on the
[flip.sol contract](https://github.com/makerdao/dss/blob/master/src/flip.sol#L117), the method signature looks like this: `tend(uint id, uint lot, uint bid)`
      - the LogNote event will take the first 




1. Persist the log record to VulcanizeDB.
  - Each event log has it's own table in the database, as well as it's
    own column in the `checked_headers` table.
    - The `checked_headers` table alllows us to keep track of which
      headers have been queried for a given log type.
  - To create a new migration file: `migrate create -ext sql -dir ./db/migrations/ create_flop_kick`
    - See `db/migrations/1536942529_create_flop_kick.up.sql`.
    - The specific log event tables are all created in the `maker`
      schema.
    - There is a one-many association between `headers` and the log
      event tables. This is so that if a header is removed due to a
reorg, the associated log event records are also removed.
  - We have been following a repository pattern to interact with the
    database for each table, see the `Create` method in `pkg/transformers/flop_kick/repository.go`.
1. Get all MissingHeaders. //TODO//
  - The repository is also responsible for querying for all header
    records that have not yet been.
1. MarkHeaderChecked//TODO//

1. Wire each component up in the transformer.

  - Each transformer's `Execute` method iterates through all of the
    MissingHeaders
  - Currently each transformer's `Execute` method is responsible for
    fetching the missing headers, meaning the headers that haven't yet
been checked for a given log type



#### For LogNote events

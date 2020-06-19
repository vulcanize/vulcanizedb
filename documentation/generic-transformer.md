# Generic Transformer
The `contractWatcher` command is a built-in generic contract watcher. It can watch events for a given contract provided the contract's ABI is available.

This command requires the contract ABI be available on Etherscan if it is not provided in the config file by the user.

## Configuration
This command takes a config of the form:

```toml
  [database]
    name     = "vulcanize_public"
    hostname = "localhost"
    port     = 5432

  [client]
    ipcPath  = "/Users/user/Library/Ethereum/geth.ipc"

  [contract]
    network  = ""
    addresses  = [
        "contractAddress1",
        "contractAddress2"
    ]
    [contract.contractAddress1]
        abi    = 'ABI for contract 1'
        startingBlock = 982463
    [contract.contractAddress2]
        abi    = 'ABI for contract 2'
        events = [
            "event1",
            "event2"
        ]
		eventArgs = [
			"arg1",
			"arg2"
		]
        startingBlock = 4448566
````

- The `contract` section defines which contracts we want to watch and with which conditions.
- `network` is only necessary if the ABIs are not provided and wish to be fetched from Etherscan.
    - Empty or nil string indicates mainnet
    - "ropsten", "kovan", and "rinkeby" indicate their respective networks
- `addresses` lists the contract addresses we are watching and is used to load their individual configuration parameters
- `contract.<contractAddress>` are the sub-mappings which contain the parameters specific to each contract address
    - `abi` is the ABI for the contract; if none is provided the application will attempt to fetch one from Etherscan using the provided address and network
    - `events` is the list of events to watch
        - If this field is omitted or no events are provided then by defualt ALL events extracted from the ABI will be watched
        - If event names are provided then only those events will be watched
    - `eventArgs` is the list of arguments to filter events with
        - If this field is omitted or no eventArgs are provided then by default watched events are not filtered by their argument values
        - If eventArgs are provided then only those events which emit at least one of these values as an argument are watched
    - `startingBlock` is the block we want to begin watching the contract, usually the deployment block of that contract

At the very minimum, for each contract address an ABI and a starting block number need to be provided (or just the starting block if the ABI can be reliably fetched from Etherscan).
With just this information we will be able to watch events on the contract.

## Output

Transformed events are committed to Postgres in schemas and tables generated according to the contract abi.

Schemas are created for each contract using the naming convention `<sync-type>_<lowercase contract-address>`.
Under this schema, tables are generated for watched events as `<lowercase event name>_event`.

## Example:

Modify `./environments/example.toml` to replace the empty `ipcPath` with a path that points to an ethjson_rpc endpoint (e.g. a local geth node ipc path or an Infura url).

If you are operating a header sync vDB, run:

 `./vulcanizedb contractWatcher --config=./environments/example.toml`

This will run the contractWatcher and configures it to watch the contracts specified in the config file.

The example config we link to in this example watches two contracts, the ENS Registry (0x314159265dD8dbb310642f98f50C066173C1259b) and TrueUSD (0x8dd5fbCe2F6a956C3022bA3663759011Dd51e73E).

Because the ENS Registry is configured with only an ABI and a starting block, we will watch all events for this contract.
Note that the ENS Registry is an example of a contract which does not have its ABI available over Etherscan and must have it included in the config file.

The TrueUSD contract is configured with two events (`Transfer` and `Mint`), as such it will watch these two events.
Note that we do not provide an ABI for TrueUSD as its ABI can be fetched from Etherscan.

For the ENS contract, it produces and populates a schema with four tables"
`header_0x314159265dd8dbb310642f98f50c066173c1259b.newowner_event`
`header_0x314159265dd8dbb310642f98f50c066173c1259b.newresolver_event`
`header_0x314159265dd8dbb310642f98f50c066173c1259b.newttl_event`
`header_0x314159265dd8dbb310642f98f50c066173c1259b.transfer_event`

For the TrusUSD contract, it produces and populates a schema with two tables:

`header_0x8dd5fbce2f6a956c3022ba3663759011dd51e73e.transfer_event`
`header_0x8dd5fbce2f6a956c3022ba3663759011dd51e73e.mint_event`

Column ids and types for these tables are generated based on the event names, resulting in tables such as:

Table "header_0x8dd5fbce2f6a956c3022ba3663759011dd51e73e.transfer_event"

|  Column    |         Type          | Collation | Nullable |                                           Default                                            | Storage  | Stats target | Description  
|:----------:|:---------------------:|:---------:|:--------:|:--------------------------------------------------------------------------------------------:|:--------:|:------------:|:-----------:|
| id         | integer               |           | not null | nextval('header_0x8dd5fbce2f6a956c3022ba3663759011dd51e73e.transfer_event_id_seq'::regclass) | plain    |              |             |
| header_id  | integer               |           | not null |                                                                                              | plain    |              |             |
| raw_log    | jsonb                 |           |          |                                                                                              | extended |              |             |
| log_idx    | integer               |           | not null |                                                                                              | plain    |              |             |
| tx_idx     | integer               |           | not null |                                                                                              | plain    |              |             |
| from_      | character varying(66) |           | not null |                                                                                              | extended |              |             |
| to_        | character varying(66) |           | not null |                                                                                              | extended |              |             |
| value_     | numeric               |           | not null |                                                                                              | main     |              |             |

The contract address used for the schema name has been down-cased.
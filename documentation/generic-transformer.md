# Generic Transformer
The `contractWatcher` command is a built-in generic contract watcher. It can watch any and all events for a given contract provided the contract's ABI is available.
It also provides some state variable coverage by automating polling of public methods, with some restrictions:
1. The method must have 2 or less arguments
1. The method's arguments must all be of type address or bytes32 (hash)
1. The method must return a single value

This command operates in two modes- `header` and `full`- which require a header or full-synced vulcanizeDB, respectively.

This command requires the contract ABI be available on Etherscan if it is not provided in the config file by the user.

If method polling is turned on we require an archival node at the ETH ipc endpoint in our config, whether or not we are operating in `header` or `full` mode.
Otherwise we only need to connect to a full node.

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
        methods = [
            "method1",
			"method2"
        ]
		methodArgs = [
			"arg1",
			"arg2"
		]
        startingBlock = 4448566
        piping = true
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
    - `methods` is the list of methods to poll
        - If this is omitted or no methods are provided then by default NO methods are polled
        - If method names are provided then those methods will be polled, provided
            1) Method has two or less arguments
            1) Arguments are all of address or hash types
            1) Method returns a single value
    - `methodArgs` is the list of arguments to limit polling methods to
        - If this field is omitted or no methodArgs are provided then by default methods will be polled with every combination of the appropriately typed values that have been collected from watched events
        - If methodArgs are provided then only those values will be used to poll methods
    - `startingBlock` is the block we want to begin watching the contract, usually the deployment block of that contract
    - `piping` is a boolean flag which indicates whether or not we want to pipe return method values forward as arguments to subsequent method calls

At the very minimum, for each contract address an ABI and a starting block number need to be provided (or just the starting block if the ABI can be reliably fetched from Etherscan).
With just this information we will be able to watch all events at the contract, but with no additional filters and no method polling.

## Output

Transformed events and polled method results are committed to Postgres in schemas and tables generated according to the contract abi.      

Schemas are created for each contract using the naming convention `<sync-type>_<lowercase contract-address>`
Under this schema, tables are generated for watched events as `<lowercase event name>_event` and for polled methods as `<lowercase method name>_method`  
The 'method' and 'event' identifiers are tacked onto the end of the table names to prevent collisions between methods and events of the same lowercase name    

## Example:

Modify `./environments/example.toml` to replace the empty `ipcPath` with a path that points to an ethjson_rpc endpoint (e.g. a local geth node ipc path or an Infura url).
This endpoint should be for an archival eth node if we want to perform method polling as this configuration is currently set up to do. To work with a non-archival full node,
remove the `balanceOf` method from the `0x8dd5fbce2f6a956c3022ba3663759011dd51e73e` (TrueUSD) contract.

If you are operating a header sync vDB, run:

 `./vulcanizedb contractWatcher --config=./environments/example.toml --mode=header`

If instead you are operating a full sync vDB and provided an archival node IPC path, run in full mode:

 `./vulcanizedb contractWatcher --config=./environments/example.toml --mode=full`

This will run the contractWatcher and configures it to watch the contracts specified in the config file. Note that
by default we operate in `header` mode but the flag is included here to demonstrate its use.

The example config we link to in this example watches two contracts, the ENS Registry (0x314159265dD8dbb310642f98f50C066173C1259b) and TrueUSD (0x8dd5fbCe2F6a956C3022bA3663759011Dd51e73E).

Because the ENS Registry is configured with only an ABI and a starting block, we will watch all events for this contract and poll none of its methods. Note that the ENS Registry is an example
of a contract which does not have its ABI available over Etherscan and must have it included in the config file.

The TrueUSD contract is configured with two events (`Transfer` and `Mint`) and a single method (`balanceOf`), as such it will watch these two events and use any addresses it collects emitted from them
to poll the `balanceOf` method with those addresses at every block. Note that we do not provide an ABI for TrueUSD as its ABI can be fetched from Etherscan.

For the ENS contract, it produces and populates a schema with four tables"
`header_0x314159265dd8dbb310642f98f50c066173c1259b.newowner_event`
`header_0x314159265dd8dbb310642f98f50c066173c1259b.newresolver_event`
`header_0x314159265dd8dbb310642f98f50c066173c1259b.newttl_event`
`header_0x314159265dd8dbb310642f98f50c066173c1259b.transfer_event`

For the TrusUSD contract, it produces and populates a schema with three tables:

`header_0x8dd5fbce2f6a956c3022ba3663759011dd51e73e.transfer_event`
`header_0x8dd5fbce2f6a956c3022ba3663759011dd51e73e.mint_event`
`header_0x8dd5fbce2f6a956c3022ba3663759011dd51e73e.balanceof_method`

Column ids and types for these tables are generated based on the event and method argument names and types and method return types, resulting in tables such as:

Table "header_0x8dd5fbce2f6a956c3022ba3663759011dd51e73e.transfer_event"

|  Column    |         Type          | Collation | Nullable |                                           Default                                            | Storage  | Stats target | Description  
|:----------:|:---------------------:|:---------:|:--------:|:--------------------------------------------------------------------------------------------:|:--------:|:------------:|:-----------:|
| id         | integer               |           | not null | nextval('header_0x8dd5fbce2f6a956c3022ba3663759011dd51e73e.transfer_event_id_seq'::regclass) | plain    |              |             |
| header_id  | integer               |           | not null |                                                                                              | plain    |              |             |
| token_name | character varying(66) |           | not null |                                                                                              | extended |              |             |
| raw_log    | jsonb                 |           |          |                                                                                              | extended |              |             |
| log_idx    | integer               |           | not null |                                                                                              | plain    |              |             |
| tx_idx     | integer               |           | not null |                                                                                              | plain    |              |             |
| from_      | character varying(66) |           | not null |                                                                                              | extended |              |             |
| to_        | character varying(66) |           | not null |                                                                                              | extended |              |             |
| value_     | numeric               |           | not null |                                                                                              | main     |              |             |


Table "header_0x8dd5fbce2f6a956c3022ba3663759011dd51e73e.balanceof_method"

|   Column   |         Type          | Collation | Nullable |                                            Default                                             | Storage  | Stats target | Description |
|:----------:|:---------------------:|:---------:|:--------:|:----------------------------------------------------------------------------------------------:|:--------:|:------------:|:-----------:|
| id         | integer               |           | not null | nextval('header_0x8dd5fbce2f6a956c3022ba3663759011dd51e73e.balanceof_method_id_seq'::regclass) | plain    |              |             |
| token_name | character varying(66) |           | not null |                                                                                                | extended |              |             |
| block      | integer               |           | not null |                                                                                                | plain    |              |             |
| who_       | character varying(66) |           | not null |                                                                                                | extended |              |             |
| returned   | numeric               |           | not null |                                                                                                | main     |              |             |
  
The addition of '_' after table names is to prevent collisions with reserved Postgres words.

Also notice that the contract address used for the schema name has been down-cased.
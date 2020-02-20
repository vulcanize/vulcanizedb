## SuperNode Subscription

A transformer can subscribe to the SueprNode service over its ipc or ws endpoints, when subscribing the transformer
specifies the chain and a set of parameters which define which subsets of that chain's data the server should feed to them.

### Ethereum data
The `streamEthSubscribe` command serves as a simple demonstration/example of subscribing to the super-node Ethereum feed, it subscribes with a set of parameters
defined in the loaded config file, and prints the streamed data to stdout. To build transformers that subscribe to and use super-node Ethereum data,
the shared/libraries/streamer can be used. 

Usage: 

`./vulcanizedb streamEthSubscribe --config=<config_file.toml>`

The config for `streamEthSubscribe` has a set of parameters to fill the [EthSubscription config structure](../../pkg/super_node/config/eth_subscription.go)

```toml
[superNode]
    [superNode.ethSubscription]
        historicalData = true
        historicalDataOnly = false
        startingBlock = 0
        endingBlock = 0
        wsPath = "ws://127.0.0.1:8080"
        [superNode.ethSubscription.headerFilter]
            off = false
            uncles = false
        [superNode.ethSubscription.txFilter]
            off = false
            src = [
                "0xde0B295669a9FD93d5F28D9Ec85E40f4cb697BAe",
            ]
            dst = [
                "0xde0B295669a9FD93d5F28D9Ec85E40f4cb697BAe",
            ]
        [superNode.ethSubscription.receiptFilter]
            off = false
            contracts = []
            topics = [
                [
                    "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
                    "0x930a61a57a70a73c2a503615b87e2e54fe5b9cdeacda518270b852296ab1a377"
                ]
            ]
        [superNode.ethSubscription.stateFilter]
            off = false
            addresses = [
               "0xde0B295669a9FD93d5F28D9Ec85E40f4cb697BAe"
           ]
           intermediateNodes = false
        [superNode.ethSubscription.storageFilter]
            off = true
            addresses = []
            storageKeys = []
            intermediateNodes = false
```

`ethSubscription.path` is used to define the SuperNode ws url OR ipc endpoint we subscribe to

`ethSubscription.historicalData` specifies whether or not the super-node should look up historical data in its cache and
send that to the subscriber, if this is set to `false` then the super-node only streams newly synced/incoming data

`ethSubscription.historicalDataOnly` will tell the super-node to only send historical data with the specified range and
not stream forward syncing data

`ethSubscription.startingBlock` is the starting block number for the range we want to receive data in

`ethSubscription.endingBlock` is the ending block number for the range we want to receive data in;
setting to 0 means there is no end/we will continue streaming indefinitely.

`ethSubscription.headerFilter` has two sub-options: `off` and `uncles`. Setting `off` to true tells the super-node to
not send any headers to the subscriber; setting `uncles` to true tells the super-node to send uncles in addition to normal headers.

`ethSubscription.txFilter` has three sub-options: `off`, `src`, and `dst`. Setting `off` to true tells the super-node to
not send any transactions to the subscriber; `src` and `dst` are string arrays which can be filled with ETH addresses we want to filter transactions for,
if they have any addresses then the super-node will only send transactions that were sent or received by the addresses contained
in `src` and `dst`, respectively.

`ethSubscription.receiptFilter` has four sub-options: `off`, `topics`, `contracts` and `matchTxs`. Setting `off` to true tells the super-node to
not send any receipts to the subscriber; `topic0s` is a string array which can be filled with event topics we want to filter for,
if it has any topics then the super-node will only send receipts that contain logs which have that topic0. Similarly, `contracts` is
a string array which can be filled with contract addresses we want to filter for, if it contains any contract addresses the super-node will
only send receipts that correspond to one of those contracts. `matchTrxs` is a bool which when set to true any receipts that correspond to filtered for
transactions will be sent by the super-node, regardless of whether or not the receipt satisfies the `topics` or `contracts` filters.

`ethSubscription.stateFilter` has three sub-options: `off`, `addresses`, and `intermediateNodes`. Setting `off` to true tells the super-node to
not send any state data to the subscriber; `addresses` is a string array which can be filled with ETH addresses we want to filter state for,
if it has any addresses then the super-node will only send state leafs (accounts) corresponding to those account addresses. By default the super-node
only sends along state leafs, if we want to receive branch and extension nodes as well `intermediateNodes` can be set to `true`.

`ethSubscription.storageFilter` has four sub-options: `off`, `addresses`, `storageKeys`, and `intermediateNodes`. Setting `off` to true tells the super-node to
not send any storage data to the subscriber; `addresses` is a string array which can be filled with ETH addresses we want to filter storage for,
if it has any addresses then the super-node will only send storage nodes from the storage tries at those state addresses. `storageKeys` is another string
array that can be filled with storage keys we want to filter storage data for. It is important to note that the storageKeys are the actual keccak256 hashes, whereas
the addresses in the `addresses` fields are the ETH addresses and not their keccak256 hashes that serve as the actual state keys. By default the super-node
only sends along storage leafs, if we want to receive branch and extension nodes as well `intermediateNodes` can be set to `true`.
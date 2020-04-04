These are the components of a VulcanizeDB Watcher:
* Data Fetcher/Streamer sources:
  * go-ethereum
  * bitcoind
  * btcd
  * IPFS
* Transformers contain:
  * converter
  * publisher
  * indexer
* Endpoints contain:
  * api
  * backend
  * filterer
  * retriever
    * ipld_server

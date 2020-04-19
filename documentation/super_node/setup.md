# VulcanizeDB Super Node Setup
Step-by-step instructions for manually setting up and running a VulcanizeDB super node.

Steps:
1. [Postgres](#postgres)
1. [Goose](#goose)
1. [IPFS](#ipfs)
1. [Blockchain](#blockchain)
1. [VulcanizeDB](#vulcanizedb)

### Postgres
A postgresDB is needed to storing all of the data in the vulcanizedb system.
Postgres is used as the backing datastore for IPFS, and is used to index the CIDs for all of the chain data stored on IPFS.
Follow the guides [here](https://wiki.postgresql.org/wiki/Detailed_installation_guides) for setting up Postgres.

Once the Postgres server is running, we will need to make a database for vulcanizedb, e.g. `vulcanize_public`.

`createdb vulcanize_public`

For running the automated tests, also create a database named `vulcanize_testing`.

`createdb vulcanize_testing`

### Goose
We use [goose](https://github.com/pressly/goose) as our migration management tool. While it is not necessary to use `goose` for manual setup, it
is required for running the automated tests.


### IPFS
We use IPFS to store IPLD objects for each type of data we extract from on chain.

To start, download and install [IPFS](https://github.com/vulcanize/go-ipfs):

`go get github.com/ipfs/go-ipfs`

`cd $GOPATH/src/github.com/ipfs/go-ipfs`

`make install`

If we want to use Postgres as our backing datastore, we need to use the vulcanize fork of go-ipfs.

Start by adding the fork and switching over to it:

`git remote add vulcanize https://github.com/vulcanize/go-ipfs.git`

`git fetch vulcanize`

`git checkout -b postgres_update vulcanize/postgres_update`

Now install this fork of ipfs, first be sure to remove any previous installation:

`make install`

Check that is installed properly by running:

`ipfs`

You should see the CLI info/help output.

And now we initialize with the `postgresds` profile.
If ipfs was previously initialized we will need to remove the old profile first.
We also need to provide env variables for the postgres connection: 

We can either set these manually, e.g.
```bash
export IPFS_PGHOST=
export IPFS_PGUSER=
export IPFS_PGDATABASE=
export IPFS_PGPORT=
export IPFS_PGPASSWORD=
```

And then run the ipfs command:

`ipfs init --profile=postgresds`

Or we can use the pre-made script at `GOPATH/src/github.com/ipfs/go-ipfs/misc/utility/ipfs_postgres.sh`
which has usage: 

`./ipfs_postgres.sh <IPFS_PGHOST> <IPFS_PGPORT> <IPFS_PGUSER> <IPFS_PGDATABASE>"`

and will ask us to enter the password, avoiding storing it to an ENV variable.

Once we have initialized ipfs, that is all we need to do with it- we do not need to run a daemon during the subsequent processes (in fact, we can't).

### Blockchain
This section describes how to setup an Ethereum or Bitcoin node to serve as a data source for the super node

#### Ethereum
For Ethereum, we currently *require* [a special fork of go-ethereum](https://github.com/vulcanize/go-ethereum/tree/statediff_at_anyblock-1.9.11). This can be setup as follows.
Skip this steps if you already have access to a node that displays the statediffing endpoints.

Begin by downloading geth and switching to the vulcanize/rpc_statediffing branch:

`go get github.com/ethereum/go-ethereum`

`cd $GOPATH/src/github.com/ethereum/go-ethereum`

`git remote add vulcanize https://github.com/vulcanize/go-ethereum.git`

`git fetch vulcanize`

`git checkout -b statediffing vulcanize/statediff_at_anyblock-1.9.11`

Now, install this fork of geth (make sure any old versions have been uninstalled/binaries removed first):

`make geth`

And run the output binary with statediffing turned on:

`cd $GOPATH/src/github.com/ethereum/go-ethereum/build/bin`

`./geth --statediff --statediff.streamblock --ws --syncmode=full`

Note: if you wish to access historical data (perform `backFill`) then the node will need to operate as an archival node (`--gcmode=archive`)

Note: other CLI options- statediff specific ones included- can be explored with `./geth help`

The output from geth should mention that it is `Starting statediff service` and block synchronization should begin shortly thereafter.
Note that until it receives a subscriber, the statediffing process does nothing but wait for one. Once a subscription is received, this
will be indicated in the output and node will begin processing and sending statediffs.

Also in the output will be the endpoints that we will use to interface with the node.
The default ws url is "127.0.0.1:8546" and the default http url is "127.0.0.1:8545".
These values will be used as the `ethereum.wsPath` and `ethereum.httpPath` in the super node config, respectively.

#### Bitcoin
For Bitcoin, the super node is able to operate entirely through the universally exposed JSON-RPC interfaces.
This means we can use any of the standard full nodes (e.g. bitcoind, btcd) as our data source.

Point at a remote node or set one up locally using the instructions for [bitcoind](https://github.com/bitcoin/bitcoin) and [btcd](https://github.com/btcsuite/btcd).

The default http url is "127.0.0.1:8332". We will use the http endpoint as both the `bitcoin.wsPath` and `bitcoin.httpPath`
(bitcoind does not support websocket endpoints, we are currently using a "subscription" wrapper around the http endpoints)

### Vulcanizedb
Finally, we can begin the vulcanizeDB process itself.

Start by downloading vulcanizedb and moving into the repo:

`go get github.com/vulcanize/vulcanizedb`

`cd $GOPATH/src/github.com/vulcanize/vulcanizedb`

Run the db migrations against the Postgres database we created for vulcanizeDB:

`goose -dir=./db/migrations postgres postgres://localhost:5432/vulcanize_public?sslmode=disable up`

At this point, if we want to run the automated tests:

`make test`
`make integration_test`

Then, build the vulcanizedb binary:

`go build`

And run the super node command with a provided [config](architecture.md/#):

`./vulcanizedb superNode --config=<config_file.toml`

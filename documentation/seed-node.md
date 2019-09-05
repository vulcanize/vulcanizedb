# Seed Node

Vulcanizedb can act as an index for Ethereum data stored on IPFS through the use of the `syncAndPublish` and
`syncPublishScreenAndServe` commands. 

## Manual Setup

These commands work in conjunction with a [state-diffing full Geth node](https://github.com/vulcanize/go-ethereum/tree/statediffing)
and IPFS.

### IPFS
To start, download and install [IPFS](https://github.com/vulcanize/go-ipfs)

`go get github.com/ipfs/go-ipfs`

`cd $GOPATH/src/github.com/ipfs/go-ipfs`

`make install`

If we want to use Postgres as our backing datastore, we need to use the vulcanize fork of go-ipfs.

Start by adding the fork and switching over to it:

`git remote add vulcanize https://github.com/vulcanize/go-ipfs.git`

`git fetch vulcanize`

`git checkout -b postgres_update vulcanize/postgres_update`

Now install this fork of ipfs, first be sure to remove any previous installation.

`make install`

Check that is installed properly by running

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

And then run the ipfs command

`ipfs init --profile=postgresds`

Or we can use the pre-made script at `GOPATH/src/github.com/ipfs/go-ipfs/misc/utility/ipfs_postgres.sh`
which has usage: 

`./ipfs_postgres.sh <IPFS_PGHOST> <IPFS_PGPORT> <IPFS_PGUSER> <IPFS_PGDATABASE>"`

and will ask us to enter the password, avoiding storing it to an ENV variable.

Once we have initialized ipfs, that is all we need to do with it- we do not need to run a daemon during the subsequent processes (in fact, we can't).

### Geth 
For Geth, we currently *require* a special fork, and we can set this up as follows:

Begin by downloading geth and switching to the vulcanize/rpc_statediffing branch

`go get github.com/ethereum/go-ethereum`

`cd $GOPATH/src/github.com/ethereum/go-ethereum`

`git remote add vulcanize https://github.com/vulcanize/go-ethereum.git`

`git fetch vulcanize`

`git checkout -b statediffing vulcanize/statediffing`

Now, install this fork of geth (make sure any old versions have been uninstalled/binaries removed first)

`make geth`

And run the output binary with statediffing turned on:

`cd $GOPATH/src/github.com/ethereum/go-ethereum/build/bin`

`./geth --statediff --statediff.streamblock --ws --syncmode=full`

Note: other CLI options- statediff specific ones included- can be explored with `./geth help`

The output from geth should mention that it is `Starting statediff service` and block synchronization should begin shortly thereafter.
Note that until it receives a subscriber, the statediffing process does essentially nothing. Once a subscription is received, this 
will be indicated in the output. 

Also in the output will be the websocket url and ipc paths that we will use to subscribe to the statediffing process.
The default ws url is "ws://127.0.0.1:8546" and the default ipcPath- on Darwin systems only- is "Users/user/Library/Ethereum/geth.ipc"

### Vulcanizedb

There are two commands to choose from:
 
#### syncAndPublish
 
`syncAndPublih` performs the functions of the seed node- syncing data from Geth, converting them to IPLDs,
publishing those IPLDs to IPFS, and creating a local Postgres index to relate their CIDS to useful metadata. 

Usage:

`./vulcanizedb syncAndPublish --config=<config_file.toml>`

The config file for the `syncAndPublish` command looks very similar to the basic config file
```toml
[database]
    name     = "vulcanize_demo"
    hostname = "localhost"
    port     = 5432

[client]
    ipcPath  = "ws://127.0.0.1:8546"
    ipfsPath = "/Users/user/.ipfs"
```

With an additional field, `client.ipcPath`, that is either the ws url or the ipc path that Geth has exposed (the url and path output
when the geth sync was started), and `client.ipfsPath` which is the path the ipfs datastore directory.

#### syncPublishScreenAndServe

`syncPublishScreenAndServe` does everything that `syncAndPublish` does, plus it opens up an RPC server which exposes
an endpoint to allow transformers to subscribe to subsets of the sync-and-published data that are relevant to their transformations

Usage:

`./vulcanizedb syncPublishScreenAndServe --config=<config_file.toml>`

The config file for the `syncPublishScreenAndServe` command has two additional fields and looks like:

```toml
[database]
    name     = "vulcanize_demo"
    hostname = "localhost"
    port     = 5432

[client]
    ipcPath  = "ws://127.0.0.1:8546"
    ipfsPath = "/Users/user/.ipfs"

[server]
    ipcPath = "/Users/user/.vulcanize/vulcanize.ipc"
    wsEndpoint = "127.0.0.1:80"
```

The additional `server.ipcPath` and `server.wsEndpoint` fields are used to set what ipc endpoint and ws url
the `syncPublishScreenAndServe` rpc server will expose itself to subscribing transformers over, respectively.
Any valid and available path and endpoint is acceptable, but keep in mind that this path and endpoint need to be known by transformers for them to subscribe to the seed node.


## Dockerfile Setup

The below provides step-by-step directions for how to setup the seed node using the provided Dockerfile on an AWS Linux AMI instance.
Note that the instance will need sufficient memory and storage for this to work.

1. Install basic dependencies 
```
sudo yum update
sudo yum install -y curl gpg gcc gcc-c++ make git
```

2. Install Go 1.12
```
wget https://dl.google.com/go/go1.12.6.linux-amd64.tar.gz
tar -xzf go1.12.6.linux-amd64.tar.gz
sudo mv go /usr/local
```

3. Edit .bash_profile to export GOPATH
```
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
```

4. Install and setup Postgres
```
sudo yum install postgresql postgresql96-server
sudo service postgresql96 initdb
sudo service postgresql96 start
sudo -u postgres createuser -s ec2-user
sudo -u postgres createdb ec2-user
sudo su postgres
psql
ALTER USER "ec2-user" WITH SUPERUSER;
/q
exit
```

4b. Edit hba_file to trust connections
```
psql
SHOW hba_file;
/q
sudo vim {PATH_TO_FILE}
```

4c. Stop and restart Postgres server to affect changes
```
sudo service postgresql96 stop
sudo service postgresql96 start
```

5. Install and start Docker (exit and re-enter ec2 instance afterwards to affect changes)
```
sudo yum install -y docker
sudo service  docker start
sudo usermod -aG docker ec2-user
```

6. Fetch the repository and switch to this working branch
```
go get github.com/vulcanize/vulcanizedb
cd $GOPATH/src/github.com/vulcanize/vulcanizedb
git checkout ipfs_concurrency
```

7. Create the db
```
createdb vulcanize_public
```

8. Build and run the Docker image
```
cd $GOPATH/src/github.com/vulcanize/vulcanizedb/dockerfiles/seed_node
docker build .
docker run --network host -e VDB_PG_CONNECT=postgres://localhost:5432/vulcanize_public?sslmode=disable {IMAGE_ID}
```


## Subscribing

A transformer can subscribe to the `syncPublishScreenAndServe` service over its ipc or ws endpoints, when subscribing the transformer
specifies which subsets of the synced data it is interested in and the server will forward only these data.

The `streamSubscribe` command serves as a simple demonstration/example of subscribing to the seed-node feed, it subscribes with a set of parameters
defined in the loaded config file, and prints the streamed data to stdout. To build transformers that subscribe to and use seed-node data,
the shared/libraries/streamer can be used. 

Usage: 

`./vulcanizedb streamSubscribe --config=<config_file.toml>`

The config for `streamSubscribe` has the `subscribe` set of parameters, for example:

```toml
[subscription]
    path = "ws://127.0.0.1:8080"
    backfill = true
    backfillOnly = false
    startingBlock = 0
    endingBlock = 0
    [subscription.headerFilter]
        off = false
        finalOnly = true
    [subscription.trxFilter]
        off = false
        src = [
            "0xde0B295669a9FD93d5F28D9Ec85E40f4cb697BAe",
        ]
        dst = [
            "0xde0B295669a9FD93d5F28D9Ec85E40f4cb697BAe",
        ]
    [subscription.receiptFilter]
        off = false
        topic0s = [
            "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
            "0x930a61a57a70a73c2a503615b87e2e54fe5b9cdeacda518270b852296ab1a377"
        ]
    [subscription.stateFilter]
        off = false
        addresses = [
           "0xde0B295669a9FD93d5F28D9Ec85E40f4cb697BAe"
       ]
       intermediateNodes = false
    [subscription.storageFilter]
        off = true
        addresses = [
            "",
            ""
        ]
        storageKeys = [
            "",
            ""
        ]
        intermediateNodes = false
```

`subscription.path` is used to define the ws url OR ipc endpoint we will subscribe to the seed-node over
(the `server.ipcPath` or `server.wsEndpoint` that the seed-node has defined in their config file).

`subscription.backfill` specifies whether or not the seed-node should look up historical data in its cache and
send that to the subscriber, if this is set to `false` then the seed-node only forwards newly synced/incoming data.

`subscription.backfillOnly` will tell the seed-node to only send historical data and not stream incoming data going forward.

`subscription.startingBlock` is the starting block number for the range we want to receive data in.

`subscription.endingBlock` is the ending block number for the range we want to receive data in;
setting to 0 means there is no end/we will continue indefinitely.

`subscription.headerFilter` has two sub-options: `off` and `finalOnly`. Setting `off` to true tells the seed-node to
not send any headers to the subscriber; setting `finalOnly` to true tells the seed-node to send only canonical headers.

`subscription.trxFilter` has three sub-options: `off`, `src`, and `dst`. Setting `off` to true tells the seed-node to
not send any transactions to the subscriber; `src` and `dst` are string arrays which can be filled with ETH addresses we want to filter transactions for,
if they have any addresses then the seed-node will only send transactions that were sent or received by the addresses contained
in `src` and `dst`, respectively.

`subscription.receiptFilter` has two sub-options: `off` and `topics`. Setting `off` to true tells the seed-node to
not send any receipts to the subscriber; `topic0s` is a string array which can be filled with event topics we want to filter for,
if it has any topics then the seed-node will only send receipts that contain logs which have that topic0.

`subscription.stateFilter` has three sub-options: `off`, `addresses`, and `intermediateNodes`. Setting `off` to true tells the seed-node to
not send any state data to the subscriber; `addresses` is a string array which can be filled with ETH addresses we want to filter state for,
if it has any addresses then the seed-node will only send state leafs (accounts) corresponding to those account addresses. By default the seed-node
only sends along state leafs, if we want to receive branch and extension nodes as well `intermediateNodes` can be set to `true`.

`subscription.storageFilter` has four sub-options: `off`, `addresses`, `storageKeys`, and `intermediateNodes`. Setting `off` to true tells the seed-node to
not send any storage data to the subscriber; `addresses` is a string array which can be filled with ETH addresses we want to filter storage for,
if it has any addresses then the seed-node will only send storage nodes from the storage tries at those state addresses. `storageKeys` is another string
array that can be filled with storage keys we want to filter storage data for. It is important to note that the storageKeys are the actual keccak256 hashes, whereas
the addresses in the `addresses` fields are the ETH addresses and not their keccak256 hashes that serve as the actual state keys. By default the seed-node
only sends along storage leafs, if we want to receive branch and extension nodes as well `intermediateNodes` can be set to `true`.

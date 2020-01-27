## Super Node Setup

Vulcanizedb can act as an index for chain data stored on IPFS through the use of the `superNode` command. 

### Manual Setup

These commands work in conjunction with a [state-diffing full Geth node](https://github.com/vulcanize/go-ethereum/tree/statediffing)
and IPFS.

#### IPFS
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

#### Geth 
For Geth, we currently *require* a special fork, and we can set this up as follows:

Begin by downloading geth and switching to the vulcanize/rpc_statediffing branch

`go get github.com/ethereum/go-ethereum`

`cd $GOPATH/src/github.com/ethereum/go-ethereum`

`git remote add vulcanize https://github.com/vulcanize/go-ethereum.git`

`git fetch vulcanize`

`git checkout -b statediffing vulcanize/statediff_at_anyblock-1.9.9`

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

#### Vulcanizedb

The `superNode` command is used to initialize and run an instance of the VulcanizeDB SuperNode

Usage:

`./vulcanizedb superNode --config=<config_file.toml`
 
 
The config file contains the parameters needed to initialize a SuperNode with the appropriate chain, settings, and services
`./vulcanizedb syncAndPublish --config=<config_file.toml>`

```toml
[superNode]
    chain = "ethereum"
    ipfsPath = "/root/.ipfs"

    [superNode.database]
        name     = "vulcanize_public"
        hostname = "localhost"
        port     = 5432
        user     = "ec2-user"

    [superNode.sync]
        on = true
        wsPath  = "ws://127.0.0.1:8546"
        workers = 1

    [superNode.server]
        on = true
        ipcPath = "/root/.vulcanize/vulcanize.ipc"
        wsPath = "127.0.0.1:8080"

    [superNode.backFill]
        on = false
        httpPath = ""
        frequency = 5
        batchSize = 50
```

### Dockerfile Setup

The below provides step-by-step directions for how to setup the super node using the provided Dockerfile on an AWS Linux AMI instance.
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
\q
exit
```

4b. Edit hba_file to trust local connections
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

6. Fetch the repository
```
go get github.com/vulcanize/vulcanizedb
cd $GOPATH/src/github.com/vulcanize/vulcanizedb
```

7. Create the db
```
createdb vulcanize_public
```

8. Build and run the Docker image
```
cd $GOPATH/src/github.com/vulcanize/vulcanizedb/dockerfiles/super_node
docker build .
docker run --network host -e IPFS_INIT=true -e VDB_PG_NAME=vulcanize_public -e VDB_PG_HOSTNAME=localhost -e VDB_PG_PORT=5432 -e VDB_PG_USER=postgres -e VDB_PG_PASSWORD=password {IMAGE_ID}
```
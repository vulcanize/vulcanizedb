#!/bin/sh
# Runs the db migrations and starts the super node services

# Exit if the variable tests fail
set -e
set +x

# Check the database variables are set
# XXX set defaults, don't silently fail
#test $DATABASE_HOSTNAME
#test $DATABASE_NAME
#test $DATABASE_PORT
#test $DATABASE_USER
#test $DATABASE_PASSWORD
#test $IPFS_INIT
#test $IPFS_PATH
VDB_COMMAND=${VDB_COMMAND:-superNode}
set +e

# Construct the connection string for postgres
VDB_PG_CONNECT=postgresql://$DATABASE_USER:$DATABASE_PASSWORD@$DATABASE_HOSTNAME:$DATABASE_PORT/$DATABASE_NAME?sslmode=disable

# Run the DB migrations
echo "Connecting with: $VDB_PG_CONNECT"
echo "Running database migrations"
./goose -dir migrations/vulcanizedb postgres "$VDB_PG_CONNECT" up
rv=$?

if [ $rv != 0 ]; then
  echo "Could not run migrations. Are the database details correct?"
  exit 1
fi

# Export our database variables so that the IPFS Postgres plugin can use them
export IPFS_PGHOST=$DATABASE_HOSTNAME
export IPFS_PGUSER=$DATABASE_USER
export IPFS_PGDATABASE=$DATABASE_NAME
export IPFS_PGPORT=$DATABASE_PORT
export IPFS_PGPASSWORD=$DATABASE_PASSWORD


if [ ! -d "$HOME/.ipfs" ]; then
  # initialize PG-IPFS
  echo "Initializing Postgres-IPFS profile"
  ./ipfs init --profile=postgresds

  rv=$?
  if [ $rv != 0 ]; then
    echo "Could not initialize ipfs"
    exit 1
  fi
fi


echo "Beginning the vulcanizedb process"
VDB_CONFIG_FILE=${VDB_CONFIG_FILE:-config.toml}
DEFAULT_OPTIONS="--config=$VDB_CONFIG_FILE"
VDB_FULL_CL=${VDB_FULL_CL:-$VDB_COMMAND $DEFAULT_OPTIONS}
echo running: ./vulcanizedb $VDB_FULL_CL $@

vdb_args="$@"
# default is to use the config passed by the build arg
if [[ -z "$vdb_args" ]];
  vdb_args="--config=config.toml"
fi

./vulcanizedb $vdb_args
rv=$?

if [ $rv != 0 ]; then
  echo "VulcanizeDB startup failed"
  exit 1
fi
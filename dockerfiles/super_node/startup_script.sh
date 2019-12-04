#!/bin/sh
# Runs the db migrations and starts the super node services

# Exit if the variable tests fail
set -e

# Check the database variables are set
test $VDB_PG_NAME
test $VDB_PG_HOSTNAME
test $VDB_PG_PORT
test $VDB_PG_USER
set +e

# Export our database variables so that the IPFS Postgres plugin can use them
export IPFS_PGHOST=$VDB_PG_HOSTNAME
export IPFS_PGUSER=$VDB_PG_USER
export IPFS_PGDATABASE=$VDB_PG_NAME
export IPFS_PGPORT=$VDB_PG_PORT
export IPFS_PGPASSWORD=$VDB_PG_PASSWORD

# Construct the connection string for postgres
VDB_PG_CONNECT=postgresql://$VDB_PG_USER:$VDB_PG_PASSWORD@$VDB_PG_HOSTNAME:$VDB_PG_PORT/$VDB_PG_NAME?sslmode=disable

# Run the DB migrations
echo "Connecting with: $VDB_PG_CONNECT"
echo "Running database migrations"
./goose -dir migrations/vulcanizedb postgres "$VDB_PG_CONNECT" up

# If the db migrations ran without err
if [ $? -eq 0 ]; then
    # Initialize PG-IPFS
    echo "Initializing Postgres-IPFS profile"
    ./ipfs init --profile=postgresds
else
    echo "Could not run migrations. Are the database details correct?"
    exit
fi

# If IPFS initialization was successful
if [ $? -eq 0 ]; then
    # Begin the state-diffing Geth process
    echo "Beginning the state-diffing Geth process"
    ./geth --statediff --statediff.streamblock --ws --syncmode=full 2>&1 | tee -a log.txt &
    sleep 1
else
    echo "Could not initialize Postgres backed IPFS profile. Are the database details correct?"
    exit
fi

# If Geth startup was successful
if [ $? -eq 0 ]; then
    # Wait until block synchronisation has begun
    echo "Waiting for block synchronization to begin"
    ( tail -f -n0 log.txt & ) | grep -q "Block synchronisation started" # this blocks til we see "Block synchronisation started"
    # And then spin up the syncPublishScreenAndServe Vulcanizedb service
    echo "Beginning the syncPublishScreenAndServe vulcanizedb process"
    ./vulcanizedb syncPublishScreenAndServe --config=config.toml 2>&1 | tee -a log.txt &
else
    echo "Could not initialize state-diffing Geth."
    exit
fi

# If Vulcanizedb startup was successful
if [ $? -eq 0 ]; then
    echo "Seed node successfully booted"
else
    echo "Could not start vulcanizedb syncPublishScreenAndServe process. Is the config file correct?"
    exit
fi

wait

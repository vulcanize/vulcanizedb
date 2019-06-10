#!/usr/bin/env bash
# Runs the migrations and starts the syncPublishScreenAndServe service

# Exit if the variable tests fail
set -e

# Check the database variables are set
test $DATABASE_NAME
test $DATABASE_HOSTNAME
test $DATABASE_PORT
test $DATABASE_USER
test $DATABASE_PASSWORD

# Export our database variables so that the IPFS Postgres plugin can use them
export IPFS_PGHOST=$DATABASE_HOSTNAME
export IPFS_PGUSER=$DATABASE_USER
export IPFS_PGDATABASE=$DATABASE_NAME
export IPFS_PGPORT=$DATABASE_PORT
export IPFS_PGPASSWORD=$DATABASE_PASSWORD

# Construct the connection string for postgres
CONNECT_STRING=postgresql://$DATABASE_USER:$DATABASE_PASSWORD@$DATABASE_HOSTNAME:$DATABASE_PORT/$DATABASE_NAME?sslmode=disable
echo "Connecting with: $CONNECT_STRING"

set +e

# Run the DB migrations
./goose postgres "$CONNECT_STRING" up
if [ $? -eq 0 ]; then
  # Fire up the services
  ipfs ipfs init --profile=postgresds
  geth --statediff --statediff.streamblock --ws --syncmode=full
  ./vulcanizedb syncPublishScreenAndServe --config environments/seedNodeStaging.toml &
else
  echo "Could not run migrations. Are the database details correct?"
fi
wait

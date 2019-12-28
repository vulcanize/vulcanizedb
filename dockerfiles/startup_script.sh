#!/bin/sh
# Runs the migrations and starts the headerSync and continuousLogSync services

# DEBUG
set -x

if test -z "$VDB_PG_CONNECT"; then
  # Exit if the variable tests fail
  set -e

  # Check the database variables are set
  test $VDB_PG_NAME
  test $VDB_PG_HOSTNAME
  test $VDB_PG_PORT
  test $VDB_PG_USER
  test $VDB_PG_PASSWORD
  set +e

  # Construct the connection string for postgres
  VDB_PG_CONNECT=postgresql://$VDB_PG_USER:$VDB_PG_PASSWORD@$VDB_PG_HOSTNAME:$VDB_PG_PORT/$VDB_PG_NAME?sslmode=disable
fi

# Run the DB migrations
echo "Connecting with: $VDB_PG_CONNECT"
./goose -dir migrations/vulcanizedb postgres "$VDB_PG_CONNECT" up

if [ $? -ne 0 ]; then
  echo "Could not run migrations. Are the database details correct?"
  exit 1
fi

# Fire up the services
for command in $VDB_COMMAND; do
  ./vulcanizedb $command --config config.toml &
done

wait

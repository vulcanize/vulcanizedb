#!/bin/sh
# Runs the migrations and starts the headerSync and continuousLogSync services

# DEBUG
set -x

if test -z "$VDB_PG_CONNECT"; then
  # Exit if the variable tests fail
  set -e

  # Check the database variables are set
  test $DATABASE_NAME
  test $DATABASE_HOSTNAME
  test $DATABASE_PORT
  test $DATABASE_USER
  test $DATABASE_PASSWORD
  set +e

  # Construct the connection string for postgres
  VDB_PG_CONNECT=postgresql://$DATABASE_USER:$DATABASE_PASSWORD@$DATABASE_HOSTNAME:$DATABASE_PORT/$DATABASE_NAME?sslmode=disable
fi

# Run the DB migrations
echo "Connecting with: $VDB_PG_CONNECT"
./goose -dir db/migrations postgres "$VDB_PG_CONNECT" up

if [ $? -ne 0 ]; then
  echo "Could not run migrations. Are the database details correct?"
  exit 1
fi

# Fire up the services
if [ $? -eq 0 ]; then
  # Fire up the services
  ./vulcanizedb headerSync --config config.toml -s 7218566 &
  ./vulcanizedb contractWatcher --config config.toml &
fi


wait

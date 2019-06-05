#!/bin/sh
# Runs the migrations and starts the headerSync and continuousLogSync services

# Exit if the variable tests fail
set -e

# Check the database variables are set
test $DATABASE_NAME
test $DATABASE_HOSTNAME
test $DATABASE_PORT
test $DATABASE_USER
test $DATABASE_PASSWORD

# Construct the connection string for postgres
CONNECT_STRING=postgresql://$DATABASE_USER:$DATABASE_PASSWORD@$DATABASE_HOSTNAME:$DATABASE_PORT/$DATABASE_NAME?sslmode=disable
echo "Connecting with: $CONNECT_STRING"

set +e

# Run the DB migrations
./goose postgres "$CONNECT_STRING" up
if [ $? -eq 0 ]; then
  # Fire up the services
  ./vulcanizedb headerSync --config environments/staging.toml &
  ./vulcanizedb continuousLogSync --config environments/staging.toml &
else
  echo "Could not run migrations. Are the database details correct?"
fi
wait

#!/bin/sh
# Runs the db migrations and starts the super node services

# Exit if the variable tests fail
set -e
set +x

# Check the database variables are set
test $DATABASE_HOSTNAME
test $DATABASE_NAME
test $DATABASE_PORT
test $DATABASE_USER
test $DATABASE_PASSWORD
set +e

# Construct the connection string for postgres
VDB_PG_CONNECT=postgresql://$DATABASE_USER:$DATABASE_PASSWORD@$DATABASE_HOSTNAME:$DATABASE_PORT/$DATABASE_NAME?sslmode=disable

# Run the DB migrations
echo "Connecting with: $VDB_PG_CONNECT"
echo "Running database migrations"
./goose -dir migrations/vulcanizedb postgres "$VDB_PG_CONNECT" up


# If the db migrations ran without err
if [[ $? -eq 0 ]]; then
    echo "Migrations ran successfully"
    exit 0
else
    echo "Could not run migrations. Are the database details correct?"
    exit 1
fi
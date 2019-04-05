# Contribution guidelines

Contributions are welcome! In addition to core contributions, developers are encouraged to build their own custom transformers which
can be run together with other custom transformers using the [composeAndExeucte](../../staging/documentation/composeAndExecute.md) command.

## Creating a new migration file
1. `make new_migration NAME=add_columnA_to_table1`
    - This will create a new timestamped migration file in `db/migrations`
1. Write the migration code in the created file, under the respective `goose` pragma
    - Goose automatically runs each migration in a transaction; don't add `BEGIN` and `COMMIT` statements.
1. Core migrations should be committed in their `goose fix`ed form.
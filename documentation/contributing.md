# Contribution guidelines

Contributions are welcome! Please open an Issues or Pull Request for any changes.

In addition to core contributions, developers are encouraged to build their own custom transformers which
can be run together with other custom transformers using the [composeAndExeucte](../../staging/documentation/composeAndExecute.md) command.

## Pull Requests
- `go fmt` is run as part of `make test` and `make integrationtest`, please make sure to check in the format changes.
- Ensure that new code is well tested, including integration testing if applicable.
- Make sure the build is passing.
- Update the README or any [documentation files](./) as necessary. If editing the Readme, please
conform to the
[standard-readme specification](https://github.com/RichardLitt/standard-readme).
- You may merge a Pull Request once you have an approval from core developer.

## Creating a new migration file
1. `make new_migration NAME=add_columnA_to_table1`
    - This will create a new timestamped migration file in `db/migrations`
1. Write the migration code in the created file, under the respective `goose` pragma
    - Goose automatically runs each migration in a transaction; don't add `BEGIN` and `COMMIT` statements.
1. Core migrations should be committed in their `goose fix`ed form. To do this, run `make version_migrations` which
converts timestamped migrations to migrations versioned by an incremented integer.

VulcanizeDB follows the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/version/1/4/code-of-conduct).

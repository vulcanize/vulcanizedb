# Contribution guidelines

Contributions are welcome! Please open an Issues or Pull Request for any changes.

In addition to core contributions, developers are encouraged to build their own custom transformers which
can be run together with other custom transformers using the [composeAndExecute](../../staging/documentation/custom-transformers.md) command.

## Pull Requests
- `go fmt` is run as part of `make test` and `make integrationtest`, please make sure to check in the format changes.
- Ensure that new code is well tested, including integration testing if applicable.
- Make sure the build is passing.
- Update the README or any [documentation files](./) as necessary. If editing the Readme, please
conform to the
[standard-readme specification](https://github.com/RichardLitt/standard-readme).
- Once a Pull Request has received two approvals it can be merged in by a core developer.

## Creating a new migration file
1. `make new_migration NAME=add_columnA_to_table1`
    - This will create a new timestamped migration file in `db/migrations`
1. Write the migration code in the created file, under the respective `goose` pragma
    - Goose automatically runs each migration in a transaction; don't add `BEGIN` and `COMMIT` statements.
1. Core migrations should be committed in their `goose fix`ed form. To do this, run `make version_migrations` which
converts timestamped migrations to migrations versioned by an incremented integer.

## Diagrams
- Diagrams were created with [draw.io](https://www.draw.io).
- To update a diagram:
  1. Go to [draw.io](https://www.draw.io).
  1. Click on *File > Open from* and choose the location of the diagram you want to update.
  1. Once open in draw.io, you may update it.
  1. Export the diagram to this repository's directory and commit it.


## Generating the Changelog
We use [github-changelog-generator](https://github.com/github-changelog-generator/github-changelog-generator) to generate release Changelogs. To be consistent with previous Changelogs, the following flags should be passed to the command:

```
--user vulcanize
--project vulcanizedb
--token {YOUR_GITHUB_TOKEN}
--no-issues
--usernames-as-github-logins
--since-tag {PREVIOUS_RELEASE_TAG}
```

For more information on why your github token is needed, and how to generate it see [https://github
.com/github-changelog-generator/github-changelog-generator#github-token](https://github.com/github-changelog-generator/github-changelog-generator#github-token).

## Code of Conduct
VulcanizeDB follows the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/version/1/4/code-of-conduct).

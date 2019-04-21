# Custom transformers
When the capabilities of the generic `contractWatcher` are not sufficient, custom transformers tailored to a specific
purpose can be leveraged.

Individual transformers can be composed together from any number of external repositories and executed as a single process using
the `compose` and `execute` commands or the `composeAndExecute` command.

## Writing custom transformers
For help with writing different types of custom transformers for the `composeAndExecute` set of commands, please see the below:

Storage Transformers
   * [Guide](../../staging/libraries/shared/factories/storage/README.md)
   * [Example](../../staging/libraries/shared/factories/storage/EXAMPLE.md)

Event Transformers
   * [Guide](../../staging/libraries/shared/factories/event/README.md)
   * [Example 1](https://github.com/vulcanize/ens_transformers/tree/master/transformers/registar)
   * [Example 2](https://github.com/vulcanize/ens_transformers/tree/master/transformers/registry)
   * [Example 3](https://github.com/vulcanize/ens_transformers/tree/master/transformers/resolver)

Contract Transformers
   * [Example 1](https://github.com/vulcanize/account_transformers)
   * [Example 2](https://github.com/vulcanize/ens_transformers/tree/master/transformers/domain_records)

## Preparing custom transformers to work as part of a plugin
To plug in an external transformer we need to:

1. Create a package that exports a variable `TransformerInitializer`, `StorageTransformerInitializer`, or `ContractTransformerInitializer` that are of type [TransformerInitializer](../staging/libraries/shared/transformer/event_transformer.go#L33)
or [StorageTransformerInitializer](../../staging/libraries/shared/transformer/storage_transformer.go#L31),
or [ContractTransformerInitializer](../../staging/libraries/shared/transformer/contract_transformer.go#L31), respectively
2. Design the transformers to work in the context of their [event](../staging/libraries/shared/watcher/event_watcher.go#L83),
[storage](../../staging/libraries/shared/watcher/storage_watcher.go#L53),
or [contract](../../staging/libraries/shared/watcher/contract_watcher.go#L68) watcher execution modes
3. Create db migrations to run against vulcanizeDB so that we can store the transformer output
    * Do not `goose fix` the transformer migrations, this is to ensure they are always ran after the core vulcanizedb migrations which are kept in their fixed form
    * Specify migration locations for each transformer in the config with the `exporter.transformer.migrations` fields
    * If the base vDB migrations occupy this path as well, they need to be in their `goose fix`ed form
    as they are [here](../../staging/db/migrations)

To update a plugin repository with changes to the core vulcanizedb repository, run `dep ensure` to update its dependencies.

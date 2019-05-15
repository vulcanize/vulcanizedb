# Watching Contract Storage

One approach VulcanizeDB takes to caching and indexing smart contracts is to ingest raw contract storage values.
Assuming that you are running an ethereum node that is writing contract storage changes to a CSV file, VulcanizeDB can parse them and persist the results to postgres.

## Assumptions

The current approach for caching smart contract storage diffs assumes that you are running a node that is writing contract storage diffs to a CSV file.
The CSV file is expected to have 5 columns: contract address, block hash, block number, storage key, storage value.

We have [a branch on vulcanize/parity-ethereum](https://github.com/vulcanize/parity-ethereum/tree/watch-storage-diffs) that enables running a node that writes storage diffs this way.

Looking forward, we would like to isolate this assumption as much as possible.
We may end up needing to read CSV data that is formatted differently, or reading data from a non-CSV source, and we do not want resulting changes to cascade throughout the codebase.

## Shared Code

VulcanizeDB has shared code for continuously reading from the CSV file written by the ethereum node and writing a parsed version of each row to postgres.

### Storage Watcher

The storage watcher is responsible for continuously delegating CSV rows to the appropriate transformer as they are being written by the ethereum node.
It maintains a mapping of contract addresses to transformers, and will ignore storage diff rows for contract addresses that do not have a corresponding transformer.

Storage watchers can be loaded with plugin storage transformers and executed using the `composeAndExecute` command.

### Storage Transformer

The storage transformer is responsible for converting raw contract storage hex values into useful data and writing them to postgres.
The storage transformer depends on contract-specific implementations of code capable of recognizing storage keys and writing the matching (decoded) storage value to disk.

```golang
func (transformer Transformer) Execute(row shared.StorageDiffRow) error {
	metadata, lookupErr := transformer.Mappings.Lookup(row.StorageKey)
	if lookupErr != nil {
		return lookupErr
	}
	value, decodeErr := shared.Decode(row, metadata)
	if decodeErr != nil {
		return decodeErr
	}
	return transformer.Repository.Create(row.BlockHeight, row.BlockHash.Hex(), metadata, value)
}
```

## Custom Code

In order to watch an additional smart contract, a developer must create three things:

1. Mappings - specify how to identify keys in the contract's storage trie.
1. Repository - specify how to persist a parsed version of the storage value matching the recognized storage key.
1. Instance - create an instance of the storage transformer that uses your mappings and repository.

### Mappings

```golang
type Mappings interface {
	Lookup(key common.Hash) (shared.StorageValueMetadata, error)
	SetDB(db *postgres.DB)
}
```

A contract-specific implementation of the mappings interface enables the storage transformer to fetch metadata associated with a storage key.

Storage metadata contains: the name of the variable matching the storage key, a raw version of any keys associated with the variable (if the variable is a mapping), and the variable's type.

```golang
type StorageValueMetadata struct {
	Name string
	Keys map[Key]string
	Type ValueType
}
```

Keys are only relevant if the variable is a mapping. For example, in the following Solidity code:

```solidity
pragma solidity ^0.4.0;

contract Contract {
  uint x;
  mapping(address => uint) y;
}
```

The metadata for variable `x` would not have any associated keys, but the metadata for a storage key associated with `y` would include the address used to specify that key's index in the mapping.

The `SetDB` function is required for the mappings to connect to the database.
A database connection may be desired when keys in a mapping variable need to be read from log events (e.g. to lookup what addresses may exist in `y`, above).

### Repository

```golang
type Repository interface {
	Create(blockNumber int, blockHash string, metadata shared.StorageValueMetadata, value interface{}) error
	SetDB(db *postgres.DB)
}
```

A contract-specific implementation of the repository interface enables the transformer to write the decoded storage value to the appropriate table in postgres.

The `Create` function is expected to recognize and persist a given storage value by the variable's name, as indicated on the row's metadata.
Note: we advise silently discarding duplicates in `Create` - as it's possible that you may read the same diff several times, and an error will trigger the storage watcher to queue that diff for later processing.

The `SetDB` function is required for the repository to connect to the database.

### Instance

```golang
type Transformer struct {
	Address    common.Address
	Mappings   storage_diffs.Mappings
	Repository storage_diffs.Repository
}
```

A new instance of the storage transformer is initialized with the contract-specific mappings and repository, as well as the contract's address.
The contract's address is included so that the watcher can query that value from the transformer in order to build up its mapping of addresses to transformers.

## Summary

To begin watching an additional smart contract, create a new mappings file for looking up storage keys on that contract, a repository for writing storage values from the contract, and initialize a new storage transformer instance with the mappings, repository, and contract address.

The new instance, wrapped in an initializer that calls `SetDB` on the mappings and repository, should be passed to the `AddTransformers` function on the storage watcher.
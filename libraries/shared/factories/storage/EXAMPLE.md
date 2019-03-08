# Storage Transformer Example

In the Storage Transformer README, we went over code that needs to be written to add a new storage transformer to VulcanizeDB.
In this document, we'll go over an example contract and discuss how one would go about watching its storage.

## Example Contract

For the purposes of this document, we'll be assuming that we're interested in watching the following contract:

```solidity
pragma solidity ^0.5.1;

contract Contract {
  uint256 public num_addresses;
  mapping(address => uint) public addresses;
  
  event AddressAdded(
    address addr,
    uint256 num_addrs
  );
  
  constructor() public {
    addresses[msg.sender] = 1;
    num_addresses = 1;
  }
  
  function add_address(address addr) public {
    bool exists = addresses[addr] > 0;
    addresses[addr]++;
    if (!exists) {
      emit AddressAdded(addr, ++num_addresses);
    }
  }
}
```

Disclaimer: this contract has not been audited and is not intended to be modeled or used in production. :)

This contract persists two values in it's storage:

1. `num_addresses`: the total number of unique addresses known to the contract.
2. `addresses`: a mapping that records the number of times an address has been added to the contract.

It also emits an event each time a new address is added into the contract's storage.

## Custom Code

In order to monitor the state of this smart contract, we'd need to implement: an event transformer, a mappings namespace, and a repository.
We will go through each of these in turn.

### Event Transformer

Given that the contract's storage includes a mapping, `addresses`, we will need to be able to identify the keys to that mapping that exist in the system so that we can recognize contract storage keys that correspond to non-zero values in that mapping.

The simplest way to be aware of keys used in a contract's mapping is to listen for contract events that emit the keys that are used in its mapping(s).
Since this contract includes an event, `AddressAdded`, that is emitted each time a new address is added to the `addresses` mapping, we will want to listen for those events and cache the adddresses that map to non-zero values.

Please see the event transformer README for detailed instructions about developing this code.
In short, it should be feasible to recognize `AddressAdded` events on the blockchain and parse them to keep a record of addresses that have been added to the system.

### Mappings

If we point an ethereum node at a blockchain hosting this contract and our node is equipped to write out storage changes happening on this contract, we will expect such changes to appear each time `add_address` (which modifies the `addresses` mapping) is called.

In order for those changes - which include raw hex versions of storage keys and storage values, to be useful for us - we need to know how to recognize and parse them.
Our mappings file should assist us with both of these tasks: the `Lookup` function should recognize raw storage keys and return known metadata about the storage value.

In order to perform this lookup, the mappings file should maintain its own mapping of known storage keys to the corresponding storage value metadata.
This internal mapping should contain the storage key for `num_addresses` as well as a storage key for each `addresses` key known to be associated with a non-zero value.

#### num_addresses

`num_addresses` is the first variable declared on the contract, and it is a simple (non-array, non-mapping) type.
Therefore, we know that its storage key is `0000000000000000000000000000000000000000000000000000000000000000`.
The storage key for non-array and non-mapping variables is (usually*) the index of the variable on the contract's storage.
If we see a storage diff being emitted from this contract with this storage key, we know that the `num_addresses` variable has been modified.

In this case, we would expect that the call `mappings.Lookup("0000000000000000000000000000000000000000000000000000000000000000")` would return metadata corresponding to the `num_addresses` variable.
This metadata would probably look something like:

```golang
shared.StorageValueMetadata{
    Name: "num_addresses",
    Keys: nil,
    Type: shared.Uint256,
}
```

<sup>*</sup> Occasionally, multiple variables may be packed into one storage slot, which complicates a direct translation of the index of the variable on the contract to its storage key.

#### addresses

`addresses` is the second variable declared on the contract, but it is a mapping.
Since it is a mapping, the storage key is more complex than `0000000000000000000000000000000000000000000000000000000000000001` (which would be the key for the variable if it were not an array or mapping).
Having a single storage slot for an entire mapping would not work, since there can be an arbitrary number of entries in a mapping, and a single storage value slot is constrained to 32 bytes.

The way that smart contract mappings are maintained in storage (in Solidity) is by creating a new storage key/value pair for each entry in the mapping, where the storage key is a hash of the occupied slot's key concatenated with the mapping's index on the contract.
Given an occupied slot's key, `k`, and a mapping's index on the contract, `i`, we can generate the storage key with the following code:

```golang
func GetMappingStorageKey(k, i string) string {
    return common.BytesToHash(crypto.Keccak256(common.FromHex(k + i))).Hex()
}
```

If we were to call the contract's `add_address` function with `0xde0B295669a9FD93d5F28D9Ec85E40f4cb697BAe`, we would expect to see an `AddressAdded` event emitted, with `0xde0B295669a9FD93d5F28D9Ec85E40f4cb697BAe` in its payload.
From that event, we would know that there exists in the contract's storage a storage key of:

```golang
GetMappingStorageKey("0xde0B295669a9FD93d5F28D9Ec85E40f4cb697BAe", "0000000000000000000000000000000000000000000000000000000000000001")
```

Executing the above code results in: `0x0f96a1133cfd5b94c329aa0526b5962bd791dbbfc481ca82f7d4a439e1e9bc40`.

Therefore, the first time `add_address` was called for this address, we would also expect to see a storage diff with a key of `0x0f96a1133cfd5b94c329aa0526b5962bd791dbbfc481ca82f7d4a439e1e9bc40` and a value of `0000000000000000000000000000000000000000000000000000000000000001`.
This would be the indication that in contract storage, the address `0xde0B295669a9FD93d5F28D9Ec85E40f4cb697BAe` maps to the value 1.

Given that we knew this address was a key in the mapping from our event transformer, we would expect a call to `mappings.Lookup("0x0f96a1133cfd5b94c329aa0526b5962bd791dbbfc481ca82f7d4a439e1e9bc40")` to return metadata corresponding to _this slot_ in the addresses mapping:

```golang
shared.StorageValueMetadata{
    Name: "addresses,
    Keys: map[Key]string{Address: "0xde0B295669a9FD93d5F28D9Ec85E40f4cb697BAe"},
    Type: shared.Uint256,
}
```

### Repository

Once we have recognized a storage diff, we can decode the storage value to the data's known type.
Since the metadata tells us that the above values are `uint256`, we can decode a value like `0000000000000000000000000000000000000000000000000000000000000001` to `1`.

The purpose of the contract-specific repository is to write that value to the database in a way that makes it useful for future queries.
Typically, the involves writing the block hash, block number, decoded value, and any keys in the metadata to a table.

The current repository interface has a generalized `Create` function that can accept any arbitrary storage row along with it's metadata.
This is deliberate, to facilitate shared use of the common storage transformer.
An implication of this decision is that the `Create` function typically includes a `switch` statement that selects which table to write to, as well as what data to include, based on the name of the variable as defined in the metadata.

An example implementation of `Create` for our example contract above might look like:

```golang
func (repository AddressStorageRepository) Create(blockNumber int, blockHash string, metadata shared.StorageValueMetadata, value interface{}) error {
    switch metadata.Name {
    case "num_addresses":
        _, err := repository.db.Exec(`INSERT INTO storage.num_addresses (block_hash, block_number, n) VALUES ($1, $2, $3)`,
            blockHash, blockNumber, value)
        return err
    case "addresses":
        _, err := repository.db.Exec(`INSERT INTO storage.addresses (block_hash, block_number, address, n) VALUES ($1, $2, $3, $4)`,
            blockHash, blockNumber, metadata.Keys[Address], value)
        return err
    default:
        panic(fmt.Sprintf("unrecognized contract storage name: %s", metadata.Name))
    }
}
```

## Summary

With our very simple address storing contract, we would be able to read it's storage diffs by implementing an event transformer, a mappings, and a repository.

The mappings would be able to lookup storage keys reflecting `num_addresses` or any slot in `addresses`, using addresses derived from watching the `AddressAdded` event for the latter.

The repository would be able to persist the value or `num_addresses` or any slot in `addresses`, using metadata returned from the mappings.

The mappings and repository could be plugged into the common storage transformer, enabling us to know the contract's state as it is changing.

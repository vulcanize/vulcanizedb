# ERC20 Transformers

## Description
The Transformers in this directory are associated with contract functions and events that conform to the [ERC20 Token interface](https://theethereum.wiki/w/index.php/ERC20_Token_Standard#The_ERC20_Token_Standard_Interface).

See `libraries/shared/TransformerREADME.md` for further information regarding the Transformer interface.

## Configuration
In addition to environment configuration mentioned in the main VulcanizeDB README, the ERC20 transformers also need to be configured with contract information for the desired token(s) to be watched. This configuration file is located at `./vulcanizedb/examples/erc20_watcher/config.go`.

## ERC20 Functions
The `everyblock` directory contains transformers that fetch data from the contract itself, via one of the standard functions.

Currently, the `totalSupply` function transformer has been implemented. This transformer will fetch the total supply for the given contract address and persist `total_supply` records in the database.

## Running the transformers
1. If running against a local node, make sure that the node has been started.
1. In a separate terminal run the following command:
`./vulcanizedb erc20 --config <config.toml>`

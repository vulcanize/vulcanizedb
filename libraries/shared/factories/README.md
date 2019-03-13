# Watching Contract Events
	
One approach VulcanizeDB takes to caching and indexing smart contracts is to watch contract events emitted in receipt logs.

With a light synced vDB we can watch events by iterating over headers retrieved from the synced `headers` table and using these headers to
fetch and verify relevant event logs from a full Ethereum node, keeping track of which headers we have checked for which events 
with our `checked_headers` table.

## Assumptions

This approach assumes you are running a vDB light sync which is itself run against a light Ethereum node,
this approach also assumes there is a full node available. 

Looking forward, we will be building fetchers that enable sourcing data from IPFS instead of an ETH node.

## Shared Code

VulcanizeDB has shared code built out for building and plugging in event transformers

### [Event Watcher (light sync)](https://github.com/vulcanize/maker-vulcanizedb/blob/staging/libraries/shared/watcher/event_watcher.go)

The event watcher is responsible for continuously fetching and delegating chunks of logs and their associated header to the appropriate transformers.

Using the `compose` or `composeAndExecute` command, event watchers can be loaded with plugin event transformers and execute over them.

### [Event Transformer](https://github.com/vulcanize/maker-vulcanizedb/blob/staging/libraries/shared/transformer/event_transformer.go)

The event transformer is responsible for converting event logs into more useful data objects and storing them in Postgres.
The event transformer is composed of converter and repository interfaces and a config struct:
```go
type Transformer struct {
	Config     transformer.TransformerConfig
	Converter  Converter
	Repository Repository
}
```

The event transformer executes over provided event logs at a given header.

In this process, the converter unpacks these logs into entities and then converts these entities 
to their final db models. These models are then written to the Postgres db by the repository.

```go
func (transformer Transformer) Execute(logs []types.Log, header core.Header, recheckHeaders constants.TransformerExecution) error {
	transformerName := transformer.Config.TransformerName
	config := transformer.Config

	if len(logs) < 1 {
		err := transformer.Repository.MarkHeaderChecked(header.Id)
		if err != nil {
			log.Printf("Error marking header as checked in %v: %v", transformerName, err)
			return err
		}
		return nil
	}

	entities, err := transformer.Converter.ToEntities(config.ContractAbi, logs)
	if err != nil {
		log.Printf("Error converting logs to entities in %v: %v", transformerName, err)
		return err
	}

	models, err := transformer.Converter.ToModels(entities)
	if err != nil {
		log.Printf("Error converting entities to models in %v: %v", transformerName, err)
		return err
	}

	err = transformer.Repository.Create(header.Id, models)
	if err != nil {
		log.Printf("Error persisting %v record: %v", transformerName, err)
		return err
	}

	return nil
}
```


## Custom Code

In order to watch events at a smart contract, for those events the developer must create:

1. Config - struct to hold configuration information (contract address, starting block, event name and signature).
1. Entity - struct to unpack the event log into.
1. Model - struct representing the final data model we want to write to Postgres.
1. Converter - an interface which can unpack event logs into our entities and convert those entities to our models. 
1. Repository - an interface to write our models to Postgres.
1. TransformerInitializer - a public variable which exports our configured transformer to be loaded as part of a plugin.
1. DB migrations - migrations to generate the Postgres schema, tables, views, function, etc that are needed to store and interface with the transformed data models.

The example event we will use looks like: 
```
event ExampleEvent(bytes32 indexed arg1, address indexed arg2, bytes32 arg3, uint256 arg4, uint256 arg5);
```

### Config

The config holds configuration variables for the event transformer, including a name for the transformer, the contract address
it is working at, the contract's ABI, the topic (e.g. event signature; topic0) that it is filtering for, and starting
and ending block numbers.

```go
type TransformerConfig struct {
	TransformerName     string
	ContractAddresses   []string
	ContractAbi         string
	Topic               string
	StartingBlockNumber int64
	EndingBlockNumber   int64 // Set -1 for indefinite transformer
}
```

### Entity

Entity field names for event arguments need to be exported and match the argument's name and type. LogIndex, 
TransactionIndex, and the Raw log are retained in order to link the data to it's source for downstream validation.

```go
type ExampleEntity struct {
	Arg1             common.Hash
	Arg2             common.Address
	Arg3             common.Hash
	Arg4             *big.Int
	Arg5             *big.Int
	LogIndex         uint
	TransactionIndex uint
	Raw              types.Log
}
```

### Model

Model fields are not constrained by the event log structure.
This allows us to rename our fields, decode or convert our log values into more useful types, and perform operations 
with or on the values before persisting the data to Postgres.

```go
type ExampleModel struct {
	EventHash         string
	UserAddress       string
	FractionSkimmed   string 
	Surplus           string 
	Deficit           string
	FinalPosition     string
	LogIndex          uint
	TransactionIndex  uint
	Raw               types.Log
}
```

### Converter

The converter needs to satisfy the interface. One for unpacking logs into the custom defined entities, and
another for converting those entities to their final db models.

```go
type Converter interface {
	ToEntities(contractAbi string, ethLog []types.Log) ([]interface{}, error)
	ToModels([]interface{}) ([]interface{}, error)
}
```

For the example event, this might look like:
```go
type ExampleConverter struct{}

func (ExampleConverter) ToEntities(contractAbi string, ethLogs []types.Log) ([]interface{}, error) {
	var entities []interface{}
	for _, ethLog := range ethLogs {
		entity := &ExampleEntity{}
		address := ethLog.Address
		abi, err := geth.ParseAbi(contractAbi)
		if err != nil {
			return nil, err
		}

		contract := bind.NewBoundContract(address, abi, nil, nil, nil)

		err = contract.UnpackLog(entity, "ExampleEvent", ethLog)
		if err != nil {
			return nil, err
		}

		entity.Raw = ethLog
		entity.LogIndex = ethLog.Index
		entity.TransactionIndex = ethLog.TxIndex

		entities = append(entities, *entity)
	}

	return entities, nil
}

func (converter ExampleConverter) ToModels(entities []interface{}) ([]interface{}, error) {
	var models []interface{}
	for _, entity := range entities {
        entity, ok := entity.(ExampleModel)
        if !ok {
            return nil, fmt.Errorf("entity of type %T, not %T", entity, ExampleModel{})
        }
        
        fractionSkimmed, err := hexutil.DecodeBig(entity.Arg3.Hex())
        if err != nil {
        	reuturn nil, err
        }
        position := new(big.Int)
        position.Sub(entity.Arg4, entity.Arg5)
        finalPosition := new(big.Int)
        if preTaxPosition.Sign() < 0 {
            finalPosition = position
        } else {
            skim := new(big.Int)
            skim.Div(position, fractionSkimmed)
            finalPosition = position.Sub(position, skim)
        }
        
        rawLog, err := json.Marshal(entity.Raw)
        if err != nil {
            return nil, err
        }
        
        model := ExampleModel{
            EventHash:        entity.Arg1.Hex(),
            UserAddress:      entity.Arg2.Hex(),
            FractionSkimmed:  fractionSkimmed.String(), 
            Surplus:          entity.Arg4.String(),
            Deficit:          entity.Arg5.String(),
            FinalPosition:    finalPosition,
            LogIndex:         entity.LogIndex,
            TransactionIndex: entity.TransactionIndex,
            Raw:              rawLog,
        }
        models = append(models, model)
	}
	return models, nil
}
```
Notice that in this example we have a bytes32 argument in the event that needs to be decoded to an integer before it can be worked with
to produce our hypothetical `FinalPosition` field. This is to highlight the fact that contracts can and sometimes do encode the
data types we want to work with into raw bytes. Writing custom transformers with these converters allows us to account for this.

### Repository

The repository needs to satisfy the interface and use the `Create` method to write the model to Postgres.
```go
type Repository interface {
	Create(headerID int64, models []interface{}) error
	MarkHeaderChecked(headerID int64) error
	MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error)
	RecheckHeaders(startingBlockNumber, endingBlockNUmber int64) ([]core.Header, error)
	SetDB(db *postgres.DB)
}
```

For the example event, this might look like:
```go
type ExampleRepository struct {
	db *postgres.DB
}

func (repository *ExampleRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository ExampleRepository) Create(headerID int64, models []interface{}) error {
	tx, dBaseErr := repository.db.Begin()
	if dBaseErr != nil {
		return dBaseErr
	}
	for _, model := range models {
		model, ok := model.(ExampleModel)
		if !ok {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback ", rollbackErr)
			}
			return fmt.Errorf("model of type %T, not %T", model, ExampleModel{})
		}
		
		_, execErr := tx.Exec(
			`INSERT into example_schema.example_event (header_id, event_hash, user_address, fraction_skimmed, surplus, deficit, final_position, log_idx, tx_idx, raw_log)
        			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
					ON CONFLICT (header_id, tx_idx, log_idx) DO UPDATE SET event_hash = $2, user_address = $3, fraction_skimmed = $4, surplus = $5, deficit = $6, final_position = $7, raw_log = $10;`,
			headerID, model.EventHash, model.UserAddress, model.FractonSkimmed, model.Surplus, model.Deficit, model.FinalPosition, model.LogIndex, model.TransactionIndex, model.Raw,
		)
		if execErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback ", rollbackErr)
			}
			return execErr
		}
	}

	checkHeaderErr := repo.MarkHeaderCheckedInTransaction(headerID, tx, "example_event_checked")
	if checkHeaderErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Error("failed to rollback ", rollbackErr)
		}
		return checkHeaderErr
	}

	return tx.Commit()
}

func (repository ExampleRepository) MarkHeaderChecked(headerID int64) error {
	return repo.MarkHeaderChecked(headerID, repository.db, "example_event_checked")
}

func (repository ExampleRepository) MissingHeaders(startingBlockNumber int64, endingBlockNumber int64) ([]core.Header, error) {
	return repo.MissingHeaders(startingBlockNumber, endingBlockNumber, repository.db,"example_event_checked")
}

func (repository ExampleRepository) RecheckHeaders(startingBlockNumber int64, endingBlockNumber int64) ([]core.Header, error) {
	return repo.RecheckHeaders(startingBlockNumber, endingBlockNumber, repository.db, "example_event_checked")
}
```

### TransformerInitializer

A transformer initializer variable needs to be exported from somewhere within the transformer repository so that the transformer can be
loaded as part of a plugin in the `compose` or `composeAndExecute` commands. It is important that this variable is named `TransformerInitializer` and
it must be of `type TransformerInitializer func(db *postgres.DB) EventTransformer`.

```go
var TransformerInitializer transformer.TransformerInitializer = factories.Transformer{
	Config:     exampleEventConfig,
	Converter:  ExampleConverter{},
	Repository: &ExampleRepository{},
}.NewTransformer
```

### DB migrations

We use `goose` as our migration management tool. Any Go data model that needs to be written to Postgres by the
repository needs a db migration for the corresponding Postgres data model.

Each contract or set of transformers being watched should define its own namespace with a db schema:
```postgresql
-- +goose Up
CREATE SCHEMA example_schema;


-- +goose Down
DROP SCHEMA example_schema;

```

For the example event and its resulting model, the table we write to would look like:
```postgresql
-- +goose Up
CREATE TABLE example_schema.example_event (
  id               SERIAL PRIMARY KEY,
  header_id        INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  event_hash       CHARACTER VARYING(66) NOT NULL,
  user_address     CHARACTER VARYING(66) NOT NULL,
  fraction_skimmed NUMERIC NOT NULL,
  surplus          NUMERIC NOT NULL,
  deficit          NUMERIC NOT NULL,
  final_position   NUMERIC NOT NULL,
  tx_idx           INTEGER NOT NUll,
  log_idx          INTEGER NOT NUll,
  raw_log          JSONB,
  UNIQUE (header_id, tx_idx, log_idx)
);

ALTER TABLE public.checked_headers
  ADD COLUMN example_event_checked INTEGER NOT NULL DEFAULT 0;


-- +goose Down
DROP TABLE maker.bite;

ALTER TABLE public.checked_headers
  DROP COLUMN bite_checked;
``` 

Notice that we have also added a column to the `checked_headers` table for this event so that we can keep track
of which headers we have already filtered through for this event.

## Summary

To create a transformer for a contract event we need to create entities for unpacking the raw log, models to represent
the final data structure, a converter to mediate this unpacking and conversion between entities to models, a repository to write
these models to Postgres, db migrations to accommodate these models in Postgres, and a TransformerInitializer to export the
configured transformer and load it as a plugin to the `compose` or `composeAndExecute` commands as described in the [main readme](https://github.com/vulcanize/maker-vulcanizedb/blob/staging/README.md#composeandexecute-configuration).

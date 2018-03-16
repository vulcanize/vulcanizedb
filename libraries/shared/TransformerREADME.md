# Transformers

## Description
Transformers must be defined in order to define what events should trigger data updates and how those are performed.

## Interface

### Initializer
Accepts DB and Blockchain from Vulcanize and returns a new transformer. E.g. for a new object "Cup":
`func NewCupTransformer(db *postgres.DB, blockchain core.ContractDataFetcher) transformers.Transformer`

### Execute
Triggers operations to take in response to a given log event.
Can persist data from logs, fetch and persist arbitrary data from outside services (e.g. contract state), or take any number of other actions. E.g.:
`func (cupTransformer *CupTransformer) Execute() error`

## Additional Requirements
Transformers must define log filters and create them so that relevant watched events can be identified and retrieved. E.g.:
```$xslt
	{
		Name:      "CupsBite",
		FromBlock: 0,
		ToBlock:   -1,
		Address:   "0x448a5065aebb8e423f0896e6c5d525c040f59af3",
		Topics:    core.Topics{"0x40cc885400000000000000000000000000000000000000000000000000000000"},
	},
```
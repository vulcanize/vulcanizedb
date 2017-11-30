package core

type Contract struct {
	Attributes ContractAttributes
}

type Blockchain interface {
	GetBlockByNumber(blockNumber int64) Block
	SubscribeToBlocks(blocks chan Block)
	StartListening()
	StopListening()
	GetContract(contractHash string) (Contract, error)
	GetContractStateAttribute(contractHash string, attributeName string) (interface{}, error)
}

package core

type Blockchain interface {
	GetBlockByNumber(blockNumber int64) Block
	SubscribeToBlocks(blocks chan Block)
	StartListening()
	StopListening()
	GetContractAttributes(contractHash string) (ContractAttributes, error)
	GetContractStateAttribute(contractHash string, attributeName string) (*string, error)
}

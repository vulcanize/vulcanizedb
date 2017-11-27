package core

type Blockchain interface {
	GetBlockByNumber(blockNumber int64) Block
	SubscribeToBlocks(blocks chan Block)
	StartListening()
	StopListening()
	GetContractAttributes(contractHash string) ([]ContractAttribute, error)
	GetContractStateAttribute(contractHash string, attributeName string) (*string, error)
}

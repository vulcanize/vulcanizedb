package core

type Blockchain interface {
	SubscribeToBlocks(blocks chan Block)
	StartListening()
}

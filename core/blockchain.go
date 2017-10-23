package core

type Blockchain interface {
	RegisterObserver(observer BlockchainObserver)
	SubscribeToEvents()
}

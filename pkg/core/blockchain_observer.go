package core

type BlockchainObserver interface {
	NotifyBlockAdded(Block)
}

package config

// Subscription config is used by a subscribing transformer to specifiy which data to receive from the seed node
type Subscription struct {
	BackFill      bool
	BackFillOnly  bool
	StartingBlock int64
	EndingBlock   int64 // set to 0 or a negative value to have no ending block
	HeaderFilter  HeaderFilter
	TrxFilter TrxFilter
	ReceiptFilter ReceiptFilter
	StateFilter StateFilter
	StorageFilter StorageFilter
}

type HeaderFilter struct {
	Off bool
	FinalOnly bool
}

type TrxFilter struct {
	Off bool
	Src []string
	Dst []string
}

type ReceiptFilter struct {
	Off     bool
	Topic0s []string
}

type StateFilter struct {
	Off               bool
	Addresses         []string // is converted to state key by taking its keccak256 hash
	IntermediateNodes bool
}

type StorageFilter struct {
	Off               bool
	Addresses         []string
	StorageKeys       []string
	IntermediateNodes bool
}
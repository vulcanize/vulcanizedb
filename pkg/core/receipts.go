package core

type Receipt struct {
	Bloom             string
	ContractAddress   string
	CumulativeGasUsed uint64
	GasUsed           uint64
	Logs              []Log
	StateRoot         string
	Status            int
	TxHash            string
}

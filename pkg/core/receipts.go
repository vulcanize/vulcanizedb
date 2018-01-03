package core

type Receipt struct {
	Bloom             string
	ContractAddress   string
	CumulativeGasUsed int64
	GasUsed           int64
	Logs              []Log
	StateRoot         string
	Status            int
	TxHash            string
}

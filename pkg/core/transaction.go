package core

type Transaction struct {
	Hash     string `db:"hash"`
	Data     string `db:"input_data"`
	Nonce    uint64 `db:"nonce"`
	To       string `db:"tx_to"`
	From     string `db:"tx_from"`
	GasLimit uint64 `db:"gaslimit"`
	GasPrice int64  `db:"gasprice"`
	Receipt
	Value string `db:"value"`
}

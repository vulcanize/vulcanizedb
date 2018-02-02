package core

type Transaction struct {
	Hash     string `db:"tx_hash"`
	Data     string `db:"tx_input_data"`
	Nonce    uint64 `db:"tx_nonce"`
	To       string `db:"tx_to"`
	From     string `db:"tx_from"`
	GasLimit int64  `db:"tx_gaslimit"`
	GasPrice int64  `db:"tx_gasprice"`
	Receipt
	Value string `db:"tx_value"`
}

package vat_init

type VatInitModel struct {
	Ilk              string
	TransactionIndex uint   `db:"tx_idx"`
	Raw              []byte `db:"raw_log"`
}

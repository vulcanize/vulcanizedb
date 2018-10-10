package vat_fold

type VatFoldModel struct {
	Ilk              string
	Urn              string
	Rate             string
	TransactionIndex uint   `db:"tx_idx"`
	Raw              []byte `db:"raw_log"`
}

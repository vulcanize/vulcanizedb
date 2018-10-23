package vat_toll

type VatTollModel struct {
	Ilk              string
	Urn              string
	Take             string
	TransactionIndex uint   `db:"tx_idx"`
	LogIndex         uint   `db:"log_idx"`
	Raw              []byte `db:"raw_log"`
}

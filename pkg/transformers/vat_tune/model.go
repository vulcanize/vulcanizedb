package vat_tune

type VatTuneModel struct {
	Ilk              string
	Urn              string
	V                string
	W                string
	Dink             string
	Dart             string
	TransactionIndex uint   `db:"tx_idx"`
	Raw              []byte `db:"raw_log"`
}

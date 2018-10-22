package vat_grab

type VatGrabModel struct {
	Ilk              string
	Urn              string
	V                string
	W                string
	Dink             string
	Dart             string
	LogIndex         uint   `db:"log_idx"`
	TransactionIndex uint   `db:"tx_idx"`
	Raw              []byte `db:"raw_log"`
}

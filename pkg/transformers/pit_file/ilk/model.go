package ilk

type PitFileIlkModel struct {
	Ilk              string
	What             string
	Data             string
	TransactionIndex uint   `db:"tx_idx"`
	Raw              []byte `db:"raw_log"`
}

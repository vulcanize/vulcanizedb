package pit_file

type PitFileModel struct {
	Ilk              string
	What             string
	Risk             string
	TransactionIndex uint   `db:"tx_idx"`
	Raw              []byte `db:"raw_log"`
}

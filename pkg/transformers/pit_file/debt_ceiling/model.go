package debt_ceiling

type PitFileDebtCeilingModel struct {
	What             string
	Data             string
	TransactionIndex uint   `db:"tx_idx"`
	Raw              []byte `db:"raw_log"`
}

package stability_fee

type PitFileStabilityFeeModel struct {
	What             string
	Data             string
	TransactionIndex uint   `db:"tx_idx"`
	Raw              []byte `db:"raw_log"`
}

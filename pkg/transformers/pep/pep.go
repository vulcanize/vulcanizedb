package pep

type Pep struct {
	BlockNumber int64  `db:"block_number"`
	HeaderID    int64  `db:"header_id"`
	UsdValue    string `db:"usd_value"`
}

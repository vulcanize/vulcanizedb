package price_feeds

import "math/big"

type PriceUpdate struct {
	BlockNumber int64  `db:"block_number"`
	HeaderID    int64  `db:"header_id"`
	UsdValue    string `db:"usd_value"`
}

func Convert(conversion string, value string, prec int) string {
	var bgflt = big.NewFloat(0.0)
	bgflt.SetString(value)
	switch conversion {
	case "ray":
		bgflt.Quo(bgflt, Ray)
	case "wad":
		bgflt.Quo(bgflt, Ether)
	}
	return bgflt.Text('g', prec)
}

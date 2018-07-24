package pep

import (
	"math/big"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

var Ether = big.NewFloat(1e18)
var Ray = big.NewFloat(1e27)

type PepTransformer struct {
	fetcher    IPepFetcher
	repository IPepRepository
}

func NewPepTransformer(chain core.BlockChain, db *postgres.DB) PepTransformer {
	fetcher := NewPepFetcher(chain)
	repository := NewPepRepository(db)
	return PepTransformer{
		fetcher:    fetcher,
		repository: repository,
	}
}

func (transformer PepTransformer) Execute(header core.Header, headerID int64) error {
	logValue, err := transformer.fetcher.FetchPepValue(header)
	if err != nil {
		if err == ErrNoMatchingLog {
			return nil
		}
		return err
	}
	pep := getPep(logValue, header, headerID)
	return transformer.repository.CreatePep(pep)
}

func getPep(logValue string, header core.Header, headerID int64) Pep {
	valueInUSD := convert("wad", logValue, 15)
	pep := Pep{
		BlockNumber: header.BlockNumber,
		HeaderID:    headerID,
		UsdValue:    valueInUSD,
	}
	return pep
}

func convert(conversion string, value string, prec int) string {
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

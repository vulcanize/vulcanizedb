package vat_init

import (
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

var VatInitConfig = shared.TransformerConfig{
	ContractAddress:     "0x239E6f0AB02713f1F8AA90ebeDeD9FC66Dc96CD6", // temporary address from Ganache deploy
	ContractAbi:         shared.VatABI,
	Topics:              []string{shared.VatInitSignature},
	StartingBlockNumber: 0,
	EndingBlockNumber:   100,
}

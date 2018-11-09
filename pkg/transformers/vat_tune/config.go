package vat_tune

import (
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
)

var VatTuneConfig = shared.TransformerConfig{
	TransformerName:     constants.VatTuneLabel,
	ContractAddresses:   []string{constants.VatContractAddress},
	ContractAbi:         constants.VatABI,
	Topic:               constants.VatTuneSignature,
	StartingBlockNumber: 0,
	EndingBlockNumber:   -1,
}

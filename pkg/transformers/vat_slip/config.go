package vat_slip

import (
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
)

var VatSlipConfig = shared.TransformerConfig{
	TransformerName:     constants.VatSlipLabel,
	ContractAddresses:   []string{constants.VatContractAddress},
	ContractAbi:         constants.VatABI,
	Topic:               constants.VatSlipSignature,
	StartingBlockNumber: 0,
	EndingBlockNumber:   -1,
}

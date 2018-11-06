package vat_slip

import "github.com/vulcanize/vulcanizedb/pkg/transformers/shared"

var VatSlipConfig = shared.TransformerConfig{
	TransformerName:     shared.VatSlipLabel,
	ContractAddresses:   []string{shared.VatContractAddress},
	ContractAbi:         shared.VatABI,
	Topic:               shared.VatSlipSignature,
	StartingBlockNumber: 0,
	EndingBlockNumber:   10000000,
}

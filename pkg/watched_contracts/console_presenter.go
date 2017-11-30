package watched_contracts

import (
	"fmt"

	"github.com/8thlight/vulcanizedb/pkg/core"
	"github.com/ethereum/go-ethereum/common"
)

func GenerateConsoleOutput(summary *ContractSummary) string {
	return fmt.Sprintf(template(),
		summary.ContractHash,
		summary.NumberOfTransactions,
		transactionToString(summary.LastTransaction),
		attributesString(summary),
	)
}

func template() string {
	return `********************Contract Summary***********************
                      HASH: %v
    NUMBER OF TRANSACTIONS: %d
          LAST TRANSACTION:
                            %s
                ATTRIBUTES:
                            %s
	`
}

func transactionToString(transaction *core.Transaction) string {
	if transaction == nil {
		return "NONE"
	} else {
		return fmt.Sprintf(`Hash: %s
                              To: %s
                            From: %s`, transaction.Hash, transaction.To, transaction.From)
	}
}

func attributesString(summary *ContractSummary) string {
	var formattedAttributes string
	for _, attribute := range summary.Attributes {
		formattedAttributes += formatAttribute(attribute.Name, summary) + "\n" + "                            "
	}
	return formattedAttributes
}

func formatAttribute(attributeName string, summary *ContractSummary) string {
	var stringResult string
	result := summary.GetStateAttribute(attributeName)
	switch t := result.(type) {
	case common.Address:
		ca := result.(common.Address)
		stringResult = fmt.Sprintf("%s: %v", attributeName, ca.Hex())
	default:
		_ = t
		stringResult = fmt.Sprintf("%s: %v", attributeName, result)
	}
	return stringResult
}

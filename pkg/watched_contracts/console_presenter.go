package watched_contracts

import (
	"fmt"

	"github.com/8thlight/vulcanizedb/pkg/core"
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
	attributes := []string{"name", "symbol"}
	var formattedAttributes string
	for _, attribute := range attributes {
		formattedAttributes += formatAttribute(attribute, summary) + "\n" + "                            "
	}
	return formattedAttributes
}

func formatAttribute(attributeName string, summary *ContractSummary) string {
	formattedAttribute := fmt.Sprintf("%s: %s", attributeName, summary.GetStateAttribute(attributeName))
	return formattedAttribute
}

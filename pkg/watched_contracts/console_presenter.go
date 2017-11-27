package watched_contracts

import (
	"fmt"

	"github.com/8thlight/vulcanizedb/pkg/core"
)

func GenerateConsoleOutput(summary *ContractSummary) string {
	return fmt.Sprintf(template(),
		summary.ContractHash,
		summary.GetStateAttribute("name"),
		summary.NumberOfTransactions,
		transactionToString(summary.LastTransaction),
	)
}

func template() string {
	return `********************Contract Summary***********************
                      HASH: %v
                      NAME: %s
    NUMBER OF TRANSACTIONS: %d
          LAST TRANSACTION:
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

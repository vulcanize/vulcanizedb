package watched_contracts

import "fmt"

func PrintReport(summary *ContractSummary) {
	fmt.Printf(`********************Contract Summary***********************

HASH: %v
`, summary.ContractHash)
}

package contract_summary_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestContractSummary(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ContractSummary Suite")
}

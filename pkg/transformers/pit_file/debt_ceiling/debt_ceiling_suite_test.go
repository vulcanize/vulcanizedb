package debt_ceiling_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDebtCeiling(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DebtCeiling Suite")
}

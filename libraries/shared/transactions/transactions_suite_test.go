package transactions_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTransactions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shared Transactions Suite")
}

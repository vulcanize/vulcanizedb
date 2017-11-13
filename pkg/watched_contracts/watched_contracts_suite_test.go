package watched_contracts_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestWatchedContracts(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "WatchedContracts Suite")
}

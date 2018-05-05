package erc20_watcher_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestErc20Watcher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Erc20Watcher Suite")
}

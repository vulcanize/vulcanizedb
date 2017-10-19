package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestVulcanizedb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Vulcanizedb Suite")
}

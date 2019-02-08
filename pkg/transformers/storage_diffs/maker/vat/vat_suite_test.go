package vat_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVat(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Vat Suite")
}

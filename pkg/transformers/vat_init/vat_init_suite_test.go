package vat_init_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVatInit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VatInit Suite")
}

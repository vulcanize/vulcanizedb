package vat_heal_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVatHeal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VatHeal Suite")
}

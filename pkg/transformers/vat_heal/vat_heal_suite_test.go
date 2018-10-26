package vat_heal_test

import (
	"io/ioutil"
	"log"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVatHeal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VatHeal Suite")
}

var _ = BeforeSuite(func() {
	log.SetOutput(ioutil.Discard)
})

package vat_tune_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"log"
)

func TestVatTune(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VatTune Suite")
}

var _ = BeforeSuite(func() {
	log.SetOutput(ioutil.Discard)
})

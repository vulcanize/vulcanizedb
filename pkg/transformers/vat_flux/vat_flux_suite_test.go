package vat_flux_test

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVatFlux(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VatFlux Suite")
}

var _ = BeforeSuite(func() {
	log.SetOutput(ioutil.Discard)
})

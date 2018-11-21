package vat_tune_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

func TestVatTune(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VatTune Suite")
}

var _ = BeforeSuite(func() {
	log.SetOutput(ioutil.Discard)
})

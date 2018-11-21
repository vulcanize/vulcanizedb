package vat_toll_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

func TestVatToll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VatToll Suite")
}

var _ = BeforeSuite(func() {
	log.SetOutput(ioutil.Discard)
})

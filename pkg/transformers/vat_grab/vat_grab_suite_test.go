package vat_grab_test

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVatGrab(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VatGrab Suite")
}

var _ = BeforeSuite(func() {
	log.SetOutput(ioutil.Discard)
})

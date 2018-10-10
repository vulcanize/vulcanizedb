package vat_grab_test

import (
	"io/ioutil"
	"log"
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

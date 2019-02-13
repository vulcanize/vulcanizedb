package vat_test

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVat(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Vat Suite")
}

var _ = BeforeSuite(func() {
	logrus.SetOutput(ioutil.Discard)
})

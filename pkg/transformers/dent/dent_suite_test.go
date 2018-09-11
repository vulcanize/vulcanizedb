package dent_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"log"
)

func TestDent(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dent Suite")
}

var _ = BeforeSuite(func() {
	log.SetOutput(ioutil.Discard)
})

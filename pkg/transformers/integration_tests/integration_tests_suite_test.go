package integration_tests

import (
	"log"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
)

var ipc string

func TestIntegrationTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IntegrationTests Suite")
}

var _ = BeforeSuite(func() {
	ipc = "http://147.75.64.249:8545" //self hosted parity kovan node
	log.SetOutput(ioutil.Discard)
})

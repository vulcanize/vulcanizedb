package every_block_test

import (
	"io/ioutil"
	"log"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEveryBlock(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EveryBlock Suite")
}

var _ = BeforeSuite(func() {
	log.SetOutput(ioutil.Discard)
})

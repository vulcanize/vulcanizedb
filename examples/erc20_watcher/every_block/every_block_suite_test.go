package every_block_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"log"
)

func TestEveryBlock(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EveryBlock Suite")
}

var _ = BeforeSuite(func() {
	log.SetOutput(ioutil.Discard)
})

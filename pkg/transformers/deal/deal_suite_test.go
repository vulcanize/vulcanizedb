package deal_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"log"
)

func TestFlipDeal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Deal Suite")
}

var _ = BeforeSuite(func() {
	log.SetOutput(ioutil.Discard)
})

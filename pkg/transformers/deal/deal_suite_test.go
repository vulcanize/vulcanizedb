package deal_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

func TestFlipDeal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Deal Suite")
}

var _ = BeforeSuite(func() {
	log.SetOutput(ioutil.Discard)
})

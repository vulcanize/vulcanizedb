package factories_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

func TestFactories(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Factories Suite")
}

var _ = BeforeSuite(func() {
	log.SetOutput(ioutil.Discard)
})

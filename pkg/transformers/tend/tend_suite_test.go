package tend_test

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTend(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tend Suite")
}

var _ = BeforeSuite(func() {
	log.SetOutput(ioutil.Discard)
})

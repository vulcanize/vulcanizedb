package flop_kick_test

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFlopKick(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FlopKick Suite")
}

var _ = BeforeSuite(func() {
	logrus.SetOutput(ioutil.Discard)
})

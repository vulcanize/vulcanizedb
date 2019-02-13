package shared_test

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestShared(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shared Suite")
}

var _ = BeforeSuite(func() {
	logrus.SetOutput(ioutil.Discard)
})

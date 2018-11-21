package dent_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

func TestDent(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dent Suite")
}

var _ = BeforeSuite(func() {
	log.SetOutput(ioutil.Discard)
})

package history_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"testing"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestHistory(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "History Suite")
}

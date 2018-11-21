package postgres_test

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestPostgres(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Postgres Suite")
}

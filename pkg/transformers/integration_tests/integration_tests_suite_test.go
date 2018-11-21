package integration_tests

import (
	log "github.com/sirupsen/logrus"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"io/ioutil"
)

var ipc string

func TestIntegrationTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IntegrationTests Suite")
}

var _ = BeforeSuite(func() {
	testConfig := viper.New()
	testConfig.SetConfigName("staging")
	testConfig.AddConfigPath("$GOPATH/src/github.com/vulcanize/vulcanizedb/environments/")
	err := testConfig.ReadInConfig()
	ipc = testConfig.GetString("client.ipcPath")
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(ioutil.Discard)
})

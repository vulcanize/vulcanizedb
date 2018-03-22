package config_test

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

var vulcanizeConfig = []byte(`
[database]
name = "dbname"
hostname = "localhost"
port = 5432

[client]
ipcPath = "IPCPATH/geth.ipc"
`)

var _ = Describe("Loading the config", func() {

	It("reads the private config using the environment", func() {
		viper.SetConfigName("config")
		viper.AddConfigPath("$GOPATH/src/github.com/vulcanize/vulcanizedb/environments/")
		Expect(viper.Get("client.ipcpath")).To(BeNil())

		testConfig := viper.New()
		testConfig.SetConfigType("toml")
		err := testConfig.ReadConfig(bytes.NewBuffer(vulcanizeConfig))
		Expect(err).To(BeNil())
		Expect(testConfig.Get("database.hostname")).To(Equal("localhost"))
		Expect(testConfig.Get("database.name")).To(Equal("dbname"))
		Expect(testConfig.Get("database.port")).To(Equal(int64(5432)))
		Expect(testConfig.Get("client.ipcpath")).To(Equal("IPCPATH/geth.ipc"))
	})

})

package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

var _ = Describe("Loading the config", func() {

	It("reads the private config using the environment", func() {
		testConfig := viper.New()
		testConfig.SetConfigName("private")
		testConfig.AddConfigPath("$GOPATH/src/github.com/vulcanize/vulcanizedb/environments/")
		err := testConfig.ReadInConfig()
		Expect(viper.Get("client.ipcpath")).To(BeNil())
		Expect(err).To(BeNil())
		Expect(testConfig.Get("database.hostname")).To(Equal("localhost"))
		Expect(testConfig.Get("database.name")).To(Equal("vulcanize_private"))
		Expect(testConfig.Get("database.port")).To(Equal(int64(5432)))
		Expect(testConfig.Get("client.ipcpath")).To(Equal("test_data_dir/geth.ipc"))
	})

})

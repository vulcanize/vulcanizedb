package config_test

import (
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	cfg "github.com/vulcanize/vulcanizedb/pkg/config"
)

var _ = Describe("Loading the config", func() {

	It("reads the private config using the environment", func() {
		privateConfig, err := cfg.NewConfig("private")

		Expect(err).To(BeNil())
		Expect(privateConfig.Database.Hostname).To(Equal("localhost"))
		Expect(privateConfig.Database.Name).To(Equal("vulcanize_private"))
		Expect(privateConfig.Database.Port).To(Equal(5432))
		expandedPath := filepath.Join(cfg.ProjectRoot(), "test_data_dir/geth.ipc")
		Expect(privateConfig.Client.IPCPath).To(Equal(expandedPath))
	})

	It("returns an error when there is no matching config file", func() {
		config, err := cfg.NewConfig("bad-config")

		Expect(config).To(Equal(cfg.Config{}))
		Expect(err).NotTo(BeNil())
	})

	It("reads the infura config using the environment", func() {
		infuraConfig, err := cfg.NewConfig("infura")

		Expect(err).To(BeNil())
		Expect(infuraConfig.Database.Hostname).To(Equal("localhost"))
		Expect(infuraConfig.Database.Name).To(Equal("vulcanize_private"))
		Expect(infuraConfig.Database.Port).To(Equal(5432))
		Expect(infuraConfig.Client.IPCPath).To(Equal("https://mainnet.infura.io/J5Vd2fRtGsw0zZ0Ov3BL"))
	})

})

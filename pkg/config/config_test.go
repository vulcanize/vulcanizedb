package config_test

import (
	"path/filepath"

	cfg "github.com/8thlight/vulcanizedb/pkg/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

})

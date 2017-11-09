package config_test

import (
	"path/filepath"

	"github.com/8thlight/vulcanizedb/pkg/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Loading the config", func() {

	It("reads the private config using the environment", func() {
		privateConfig, err := config.NewConfig("private")

		Expect(err).To(BeNil())
		Expect(privateConfig.Database.Hostname).To(Equal("localhost"))
		Expect(privateConfig.Database.Name).To(Equal("vulcanize_private"))
		Expect(privateConfig.Database.Port).To(Equal(5432))
		expandedPath := filepath.Join(config.ProjectRoot(), "test_data_dir/geth.ipc")
		Expect(privateConfig.Client.IPCPath).To(Equal(expandedPath))
	})

	It("returns an error when there is no matching config file", func() {
		config, err := config.NewConfig("bad-config")

		Expect(config).To(BeNil())
		Expect(err).NotTo(BeNil())
	})

})

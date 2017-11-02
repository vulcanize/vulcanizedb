package config_test

import (
	"path/filepath"

	"github.com/8thlight/vulcanizedb/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Loading the config", func() {

	It("reads the private config using the environment", func() {
		privateConfig := config.NewConfig("private")

		Expect(privateConfig.Database.Hostname).To(Equal("localhost"))
		Expect(privateConfig.Database.Name).To(Equal("vulcanize_private"))
		Expect(privateConfig.Database.Port).To(Equal(5432))
		expandedPath := filepath.Join(config.ProjectRoot(), "test_data_dir/geth.ipc")
		Expect(privateConfig.Client.IPCPath).To(Equal(expandedPath))
	})

})

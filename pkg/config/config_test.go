// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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

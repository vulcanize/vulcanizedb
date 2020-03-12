// VulcanizeDB
// Copyright © 2019 Vulcanize

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

package resync

import (
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

// Config holds the parameters needed to perform a resync
type Config struct {
	Chain         shared.ChainType // The type of resync to perform
	ResyncType    shared.DataType  // The type of data to resync
	ClearOldCache bool             // Resync will first clear all the data within the range

	// DB info
	DB       *postgres.DB
	DBConfig config.Database
	IPFSPath string

	HTTPClient interface{} // Note this client is expected to support the retrieval of the specified data type(s)
	Ranges     [][2]uint64 // The block height ranges to resync
	BatchSize  uint64      // BatchSize for the resync http calls (client has to support batch sizing)

	Quit chan bool // Channel for shutting down
}

// NewReSyncConfig fills and returns a resync config from toml parameters
func NewReSyncConfig() (Config, error) {
	panic("implement me")
}

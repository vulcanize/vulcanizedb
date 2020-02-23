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

package watcher

import (
	"fmt"

	"github.com/vulcanize/vulcanizedb/libraries/shared/streamer"
	"github.com/vulcanize/vulcanizedb/pkg/eth/core"
	"github.com/vulcanize/vulcanizedb/pkg/postgres"
	shared2 "github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
	"github.com/vulcanize/vulcanizedb/pkg/watcher/eth"
	"github.com/vulcanize/vulcanizedb/pkg/watcher/shared"
)

// NewSuperNodeStreamer returns a new shared.SuperNodeStreamer
func NewSuperNodeStreamer(client core.RPCClient) shared.SuperNodeStreamer {
	return streamer.NewSuperNodeStreamer(client)
}

// NewRepository constructs and returns a new Repository that satisfies the shared.Repository interface for the specified chain
func NewRepository(chain shared2.ChainType, db *postgres.DB, triggerFuncs [][2]string) (shared.Repository, error) {
	switch chain {
	case shared2.Ethereum:
		return eth.NewRepository(db, triggerFuncs), nil
	default:
		return nil, fmt.Errorf("NewRepository constructor unexpected chain type %s", chain.String())
	}
}

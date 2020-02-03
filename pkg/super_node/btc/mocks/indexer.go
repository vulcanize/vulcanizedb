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

package mocks

import (
	"fmt"

	"github.com/vulcanize/vulcanizedb/pkg/super_node/btc"

	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"

	"github.com/vulcanize/vulcanizedb/pkg/super_node/eth"
)

// CIDIndexer is the underlying struct for the Indexer interface
type CIDIndexer struct {
	PassedCIDPayload []*btc.CIDPayload
	ReturnErr        error
}

// Index indexes a cidPayload in Postgres
func (repo *CIDIndexer) Index(cids shared.CIDsForIndexing) error {
	cidPayload, ok := cids.(*btc.CIDPayload)
	if !ok {
		return fmt.Errorf("index expected cids type %T got %T", &eth.CIDPayload{}, cids)
	}
	repo.PassedCIDPayload = append(repo.PassedCIDPayload, cidPayload)
	return repo.ReturnErr
}

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

package ipfs

import (
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/statediff"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type Syncer interface {
	Sync(payloadChan chan statediff.Payload) (*rpc.ClientSubscription, error)
}

type StateDiffSyncer struct {
	Client      core.RpcClient
	PayloadChan chan statediff.Payload
}

func NewStateDiffSyncer(client core.RpcClient) *StateDiffSyncer {
	return &StateDiffSyncer{
		Client: client,
	}
}

// Sync is the main loop for subscribing to data from the Geth state diff process
func (i *StateDiffSyncer) Sync(payloadChan chan statediff.Payload) (*rpc.ClientSubscription, error) {
	return i.Client.Subscribe("statediff", i.PayloadChan)
}
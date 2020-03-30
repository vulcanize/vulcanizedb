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

package statediff

import (
	"context"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/statediff"

	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

// APIName is the namespace for the super node's eth api
const APIName = "statediff"

// APIVersion is the version of the super node's eth api
const APIVersion = "0.0.1"

type PublicStateDiffAPI struct {
	sn shared.SuperNode
}

// NewPublicStateDiffAPI creates and returns a new PublicStateDiffAPI
func NewPublicStateDiffAPI(sn shared.SuperNode) *PublicStateDiffAPI {
	return &PublicStateDiffAPI{
		sn: sn,
	}
}

// Stream is the public method to setup a subscription that fires off statediff service payloads as they are created
func (api *PublicStateDiffAPI) Stream(ctx context.Context) (*rpc.Subscription, error) {

	return nil, nil
}

// StateDiffAt returns a statediff payload at the specific blockheight
func (api *PublicStateDiffAPI) StateDiffAt(ctx context.Context, blockNumber uint64) (*statediff.Payload, error) {
	return nil, nil
}
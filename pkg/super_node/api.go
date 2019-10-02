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

package super_node

import (
	"context"

	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/libraries/shared/streamer"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

// APIName is the namespace used for the state diffing service API
const APIName = "vdb"

// APIVersion is the version of the state diffing service API
const APIVersion = "0.0.1"

// PublicSuperNodeAPI is the public api for the super node
type PublicSuperNodeAPI struct {
	sni NodeInterface
}

// NewPublicSuperNodeAPI creates a new PublicSuperNodeAPI with the provided underlying SyncPublishScreenAndServe process
func NewPublicSuperNodeAPI(superNodeInterface NodeInterface) *PublicSuperNodeAPI {
	return &PublicSuperNodeAPI{
		sni: superNodeInterface,
	}
}

// Stream is the public method to setup a subscription that fires off SyncPublishScreenAndServe payloads as they are created
func (api *PublicSuperNodeAPI) Stream(ctx context.Context, streamFilters config.Subscription) (*rpc.Subscription, error) {
	// ensure that the RPC connection supports subscriptions
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return nil, rpc.ErrNotificationsUnsupported
	}

	// create subscription and start waiting for statediff events
	rpcSub := notifier.CreateSubscription()

	go func() {
		// subscribe to events from the SyncPublishScreenAndServe service
		payloadChannel := make(chan streamer.SuperNodePayload, payloadChanBufferSize)
		quitChan := make(chan bool, 1)
		go api.sni.Subscribe(rpcSub.ID, payloadChannel, quitChan, streamFilters)

		// loop and await state diff payloads and relay them to the subscriber with then notifier
		for {
			select {
			case packet := <-payloadChannel:
				if notifyErr := notifier.Notify(rpcSub.ID, packet); notifyErr != nil {
					log.Error("Failed to send state diff packet", "err", notifyErr)
					api.sni.Unsubscribe(rpcSub.ID)
					return
				}
			case <-rpcSub.Err():
				api.sni.Unsubscribe(rpcSub.ID)
				return
			case <-quitChan:
				// don't need to unsubscribe, SyncPublishScreenAndServe service does so before sending the quit signal
				return
			}
		}
	}()

	return rpcSub, nil
}

// Node is a public rpc method to allow transformers to fetch the Geth node info for the super node
func (api *PublicSuperNodeAPI) Node() core.Node {
	return api.sni.Node()
}

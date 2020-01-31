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

	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"

	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/eth/core"
)

// APIName is the namespace used for the state diffing service API
const APIName = "vdb"

// APIVersion is the version of the state diffing service API
const APIVersion = "0.0.1"

// PublicSuperNodeAPI is the public api for the super node
type PublicSuperNodeAPI struct {
	sn SuperNode
}

// NewPublicSuperNodeAPI creates a new PublicSuperNodeAPI with the provided underlying SyncPublishScreenAndServe process
func NewPublicSuperNodeAPI(superNodeInterface SuperNode) *PublicSuperNodeAPI {
	return &PublicSuperNodeAPI{
		sn: superNodeInterface,
	}
}

// Stream is the public method to setup a subscription that fires off super node payloads as they are processed
func (api *PublicSuperNodeAPI) Stream(ctx context.Context, params shared.SubscriptionSettings) (*rpc.Subscription, error) {
	// ensure that the RPC connection supports subscriptions
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return nil, rpc.ErrNotificationsUnsupported
	}

	// create subscription and start waiting for stream events
	rpcSub := notifier.CreateSubscription()

	go func() {
		// subscribe to events from the SyncPublishScreenAndServe service
		payloadChannel := make(chan SubscriptionPayload, PayloadChanBufferSize)
		quitChan := make(chan bool, 1)
		go api.sn.Subscribe(rpcSub.ID, payloadChannel, quitChan, params)

		// loop and await payloads and relay them to the subscriber using notifier
		for {
			select {
			case packet := <-payloadChannel:
				if err := notifier.Notify(rpcSub.ID, packet); err != nil {
					log.Error("Failed to send super node packet", "err", err)
					api.sn.Unsubscribe(rpcSub.ID)
					return
				}
			case <-rpcSub.Err():
				api.sn.Unsubscribe(rpcSub.ID)
				return
			case <-quitChan:
				// don't need to unsubscribe to super node, the service does so before sending the quit signal this way
				return
			}
		}
	}()

	return rpcSub, nil
}

// Node is a public rpc method to allow transformers to fetch the node info for the super node
func (api *PublicSuperNodeAPI) Node() core.Node {
	return api.sn.Node()
}

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

	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/eth/core"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/btc"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/eth"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
	v "github.com/vulcanize/vulcanizedb/version"
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
func (api *PublicSuperNodeAPI) Stream(ctx context.Context, rlpParams []byte) (*rpc.Subscription, error) {
	var params shared.SubscriptionSettings
	switch api.sn.Chain() {
	case shared.Ethereum:
		var ethParams eth.SubscriptionSettings
		if err := rlp.DecodeBytes(rlpParams, &ethParams); err != nil {
			return nil, err
		}
		params = &ethParams
	case shared.Bitcoin:
		var btcParams btc.SubscriptionSettings
		if err := rlp.DecodeBytes(rlpParams, &btcParams); err != nil {
			return nil, err
		}
		params = &btcParams
	default:
		panic("SuperNode is not configured for a specific chain type")
	}
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
// NOTE: this is the node info for the node that the super node is syncing from, not the node info for the super node itself
func (api *PublicSuperNodeAPI) Node() *core.Node {
	return api.sn.Node()
}

// Chain returns the chain type that this super node instance supports
func (api *PublicSuperNodeAPI) Chain() shared.ChainType {
	return api.sn.Chain()
}

// Struct for holding super node meta data
type InfoAPI struct{}

// NewPublicSuperNodeAPI creates a new PublicSuperNodeAPI with the provided underlying SyncPublishScreenAndServe process
func NewInfoAPI() *InfoAPI {
	return &InfoAPI{}
}

// Modules returns modules supported by this api
func (iapi *InfoAPI) Modules() map[string]string {
	return map[string]string{
		"vdb": "Stream",
	}
}

// NodeInfo gathers and returns a collection of metadata for the super node
func (iapi *InfoAPI) NodeInfo() *p2p.NodeInfo {
	return &p2p.NodeInfo{
		// TODO: formalize this
		ID:   "vulcanizeDB",
		Name: "superNode",
	}
}

// Version returns the version of the super node
func (iapi *InfoAPI) Version() string {
	return v.VersionWithMeta
}

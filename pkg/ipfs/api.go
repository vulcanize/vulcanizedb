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
	"context"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
)

// APIName is the namespace used for the state diffing service API
const APIName = "vulcanizedb"

// APIVersion is the version of the state diffing service API
const APIVersion = "0.0.1"

// PublicSeedNodeAPI is the public api for the seed node
type PublicSeedNodeAPI struct {
	snp SyncPublishScreenAndServe
}

// NewPublicSeedNodeAPI creates a new PublicSeedNodeAPI with the provided underlying SyncPublishScreenAndServe process
func NewPublicSeedNodeAPI(snp SyncPublishScreenAndServe) *PublicSeedNodeAPI {
	return &PublicSeedNodeAPI{
		snp: snp,
	}
}

// Subscribe is the public method to setup a subscription that fires off state-diff payloads as they are created
func (api *PublicSeedNodeAPI) Subscribe(ctx context.Context, payloadChanForTypeDefOnly chan ResponsePayload) (*rpc.Subscription, error) {
	// ensure that the RPC connection supports subscriptions
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return nil, rpc.ErrNotificationsUnsupported
	}

	streamFilters := StreamFilters{}
	streamFilters.HeaderFilter.FinalOnly = true
	streamFilters.TrxFilter.Src = []string{"0xde0B295669a9FD93d5F28D9Ec85E40f4cb697BAe"}
	streamFilters.TrxFilter.Dst = []string{"0xde0B295669a9FD93d5F28D9Ec85E40f4cb697BAe"}
	streamFilters.ReceiptFilter.Topic0s = []string{
		"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
		"0x930a61a57a70a73c2a503615b87e2e54fe5b9cdeacda518270b852296ab1a377",
	}
	streamFilters.StateFilter.Addresses = []string{"0xde0B295669a9FD93d5F28D9Ec85E40f4cb697BAe"}
	streamFilters.StorageFilter.Off = true
	// create subscription and start waiting for statediff events
	rpcSub := notifier.CreateSubscription()

	go func() {
		// subscribe to events from the state diff service
		payloadChannel := make(chan ResponsePayload)
		quitChan := make(chan bool)
		go api.snp.Subscribe(rpcSub.ID, payloadChannel, quitChan, &streamFilters)

		// loop and await state diff payloads and relay them to the subscriber with then notifier
		for {
			select {
			case packet := <-payloadChannel:
				if err := notifier.Notify(rpcSub.ID, packet); err != nil {
					log.Error("Failed to send state diff packet", "err", err)
				}
			case <-rpcSub.Err():
				err := api.snp.Unsubscribe(rpcSub.ID)
				if err != nil {
					log.Error("Failed to unsubscribe from the state diff service", err)
				}
				return
			case <-quitChan:
				// don't need to unsubscribe, statediff service does so before sending the quit signal
				return
			}
		}
	}()

	return rpcSub, nil
}
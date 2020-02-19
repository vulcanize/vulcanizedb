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

package btc

import (
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

const (
	PayloadChanBufferSize = 20000 // the max eth sub buffer size
)

// PayloadStreamer satisfies the PayloadStreamer interface for bitcoin
type PayloadStreamer struct {
	Config *rpcclient.ConnConfig
}

// NewPayloadStreamer creates a pointer to a new PayloadStreamer which satisfies the PayloadStreamer interface for bitcoin
func NewPayloadStreamer(clientConfig *rpcclient.ConnConfig) *PayloadStreamer {
	return &PayloadStreamer{
		Config: clientConfig,
	}
}

// Stream is the main loop for subscribing to data from the btc block notifications
// Satisfies the shared.PayloadStreamer interface
func (ps *PayloadStreamer) Stream(payloadChan chan shared.RawChainData) (shared.ClientSubscription, error) {
	logrus.Info("streaming block payloads from btc")
	blockNotificationHandler := rpcclient.NotificationHandlers{
		// Notification handler for block connections, forwards new block data to the payloadChan
		OnFilteredBlockConnected: func(height int32, header *wire.BlockHeader, txs []*btcutil.Tx) {
			payloadChan <- BlockPayload{
				BlockHeight: int64(height),
				Header:      header,
				Txs:         txs,
			}
		},
	}
	// Create a new client, and connect to btc ws server
	client, err := rpcclient.New(ps.Config, &blockNotificationHandler)
	if err != nil {
		return nil, err
	}
	// Register for block connect notifications.
	if err := client.NotifyBlocks(); err != nil {
		return nil, err
	}
	client.WaitForShutdown()
	return &ClientSubscription{client: client}, nil
}

// ClientSubscription is a wrapper around the underlying btcd rpc client
// to fit the shared.ClientSubscription interface
type ClientSubscription struct {
	client *rpcclient.Client
}

// Unsubscribe satisfies the rpc.Subscription interface
func (bcs *ClientSubscription) Unsubscribe() {
	bcs.client.Shutdown()
}

// Err() satisfies the rpc.Subscription interface with a dummy err channel
func (bcs *ClientSubscription) Err() <-chan error {
	errChan := make(chan error)
	return errChan
}

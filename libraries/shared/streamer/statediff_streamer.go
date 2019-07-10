// Copyright 2019 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package streamer

import (
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/statediff"
	"github.com/sirupsen/logrus"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type Streamer interface {
	Stream(chan statediff.Payload) (*rpc.ClientSubscription, error)
}

type StateDiffStreamer struct {
	client core.RpcClient
}

func (streamer *StateDiffStreamer) Stream(payloadChan chan statediff.Payload) (*rpc.ClientSubscription, error) {
	logrus.Info("streaming diffs from geth")
	return streamer.client.Subscribe("statediff", payloadChan, "stream")
}

func NewStateDiffStreamer(client core.RpcClient) StateDiffStreamer {
	return StateDiffStreamer{
		client: client,
	}
}

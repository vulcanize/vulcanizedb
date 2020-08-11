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
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/sirupsen/logrus"
)

type Streamer interface {
	Stream(chan filters.Payload) (core.Subscription, error)
}

type EthStateChangeStreamer struct {
	ethClient   core.EthClient
	filterQuery ethereum.FilterQuery
}

func NewEthStateChangeStreamer(ethClient core.EthClient, filterQuery ethereum.FilterQuery) EthStateChangeStreamer {
	return EthStateChangeStreamer{
		ethClient:   ethClient,
		filterQuery: filterQuery,
	}
}

func (streamer *EthStateChangeStreamer) Stream(payloadChan chan filters.Payload) (core.Subscription, error) {
	logrus.Info("streaming diffs from geth")
	return streamer.ethClient.SubscribeNewStateChanges(context.Background(), streamer.filterQuery, payloadChan)
}

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
	"context"

	"github.com/ethereum/go-ethereum/rpc"
)

type StreamClient struct {
	passedContext       context.Context
	passedResult        interface{}
	passedNamespace     string
	passedPayloadChan   interface{}
	passedSubscribeArgs []interface{}
}

func (client *StreamClient) Subscribe(ctx context.Context, namespace string, payloadChan interface{}, args ...interface{}) (*rpc.ClientSubscription, error) {
	client.passedNamespace = namespace
	client.passedPayloadChan = payloadChan
	client.passedContext = ctx

	for _, arg := range args {
		client.passedSubscribeArgs = append(client.passedSubscribeArgs, arg)
	}

	subscription := rpc.ClientSubscription{}
	return &subscription, nil
}

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
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/statediff"
	"github.com/vulcanize/vulcanizedb/pkg/eth/client"
)

// BackFillerClient is a mock client for use in backfiller tests
type BackFillerClient struct {
	MappedStateDiffAt map[uint64][]byte
}

// SetReturnDiffAt method to set what statediffs the mock client returns
func (mc *BackFillerClient) SetReturnDiffAt(height uint64, diffPayload statediff.Payload) error {
	if mc.MappedStateDiffAt == nil {
		mc.MappedStateDiffAt = make(map[uint64][]byte)
	}
	by, err := json.Marshal(diffPayload)
	if err != nil {
		return err
	}
	mc.MappedStateDiffAt[height] = by
	return nil
}

// BatchCall mockClient method to simulate batch call to geth
func (mc *BackFillerClient) BatchCall(batch []client.BatchElem) error {
	if mc.MappedStateDiffAt == nil {
		return errors.New("mockclient needs to be initialized with statediff payloads and errors")
	}
	for _, batchElem := range batch {
		if len(batchElem.Args) != 1 {
			return errors.New("expected batch elem to contain single argument")
		}
		blockHeight, ok := batchElem.Args[0].(uint64)
		if !ok {
			return errors.New("expected batch elem argument to be a uint64")
		}
		err := json.Unmarshal(mc.MappedStateDiffAt[blockHeight], batchElem.Result)
		if err != nil {
			return err
		}
	}
	return nil
}

// BatchCallContext mockClient method to simulate batch call to geth
func (mc *BackFillerClient) BatchCallContext(ctx context.Context, batch []client.BatchElem) error {
	if mc.MappedStateDiffAt == nil {
		return errors.New("mockclient needs to be initialized with statediff payloads and errors")
	}
	for _, batchElem := range batch {
		if len(batchElem.Args) != 1 {
			return errors.New("expected batch elem to contain single argument")
		}
		blockHeight, ok := batchElem.Args[0].(uint64)
		if !ok {
			return errors.New("expected batch elem argument to be a uint64")
		}
		err := json.Unmarshal(mc.MappedStateDiffAt[blockHeight], batchElem.Result)
		if err != nil {
			return err
		}
	}
	return nil
}

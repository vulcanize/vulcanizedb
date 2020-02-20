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
	"fmt"

	"github.com/vulcanize/vulcanizedb/pkg/super_node/btc"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

// PayloadConverter is the underlying struct for the Converter interface
type PayloadConverter struct {
	PassedStatediffPayload btc.BlockPayload
	ReturnIPLDPayload      btc.ConvertedPayload
	ReturnErr              error
}

// Convert method is used to convert a geth statediff.Payload to a IPLDPayload
func (pc *PayloadConverter) Convert(payload shared.RawChainData) (shared.ConvertedData, error) {
	stateDiffPayload, ok := payload.(btc.BlockPayload)
	if !ok {
		return nil, fmt.Errorf("convert expected payload type %T got %T", btc.BlockPayload{}, payload)
	}
	pc.PassedStatediffPayload = stateDiffPayload
	return pc.ReturnIPLDPayload, pc.ReturnErr
}

// IterativePayloadConverter is the underlying struct for the Converter interface
type IterativePayloadConverter struct {
	PassedStatediffPayload []btc.BlockPayload
	ReturnIPLDPayload      []btc.ConvertedPayload
	ReturnErr              error
	iteration              int
}

// Convert method is used to convert a geth statediff.Payload to a IPLDPayload
func (pc *IterativePayloadConverter) Convert(payload shared.RawChainData) (shared.ConvertedData, error) {
	stateDiffPayload, ok := payload.(btc.BlockPayload)
	if !ok {
		return nil, fmt.Errorf("convert expected payload type %T got %T", btc.BlockPayload{}, payload)
	}
	pc.PassedStatediffPayload = append(pc.PassedStatediffPayload, stateDiffPayload)
	if len(pc.PassedStatediffPayload) < pc.iteration+1 {
		return nil, fmt.Errorf("IterativePayloadConverter does not have a payload to return at iteration %d", pc.iteration)
	}
	returnPayload := pc.ReturnIPLDPayload[pc.iteration]
	pc.iteration++
	return returnPayload, pc.ReturnErr
}

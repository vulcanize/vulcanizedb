// VulcanizeDB
// Copyright © 2018 Vulcanize

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

package converter

import (
	"errors"

	geth "github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/contract"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/types"
)

type Converter interface {
	Convert(log geth.Log, event types.Event) (*types.Log, error)
	Update(info *contract.Contract)
}

type converter struct {
	ContractInfo *contract.Contract
}

func NewConverter(info *contract.Contract) *converter {
	return &converter{
		ContractInfo: info,
	}
}

func (c *converter) Update(info *contract.Contract) {
	c.ContractInfo = info
}

// Convert the given watched event log into a types.Log for the given event
func (c *converter) Convert(log geth.Log, event types.Event) (*types.Log, error) {
	return nil, errors.New("implement me")
}

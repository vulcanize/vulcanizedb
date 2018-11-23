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

package generic

import (
	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/constants"
)

var DaiConfig = shared.ContractConfig{
	Address:    constants.DaiContractAddress,
	Abi:        constants.DaiAbiString,
	FirstBlock: int64(4752008),
	LastBlock:  -1,
	Name:       "Dai",
	Filters:    constants.DaiERC20Filters,
}

var TusdConfig = shared.ContractConfig{
	Address:    constants.TusdContractAddress,
	Owner:      constants.TusdContractOwner,
	Abi:        constants.TusdAbiString,
	FirstBlock: int64(5197514),
	LastBlock:  -1,
	Name:       "Tusd",
	Filters:    constants.TusdGenericFilters,
}

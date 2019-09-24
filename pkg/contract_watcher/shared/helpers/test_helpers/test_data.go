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

package test_helpers

import (
	"strings"

	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/constants"
)

var ens = strings.ToLower(constants.EnsContractAddress)
var tusd = strings.ToLower(constants.TusdContractAddress)

var TusdConfig = config.ContractConfig{
	Network: "",
	Addresses: map[string]bool{
		tusd: true,
	},
	Abis: map[string]string{
		tusd: "",
	},
	Events: map[string][]string{
		tusd: {"Transfer"},
	},
	Methods: map[string][]string{
		tusd: nil,
	},
	MethodArgs: map[string][]string{
		tusd: nil,
	},
	EventArgs: map[string][]string{
		tusd: nil,
	},
	StartingBlocks: map[string]int64{
		tusd: 5197514,
	},
}

var ENSConfig = config.ContractConfig{
	Network: "",
	Addresses: map[string]bool{
		ens: true,
	},
	Abis: map[string]string{
		ens: "",
	},
	Events: map[string][]string{
		ens: {"NewOwner"},
	},
	Methods: map[string][]string{
		ens: nil,
	},
	MethodArgs: map[string][]string{
		ens: nil,
	},
	EventArgs: map[string][]string{
		ens: nil,
	},
	StartingBlocks: map[string]int64{
		ens: 3327417,
	},
}

var ENSandTusdConfig = config.ContractConfig{
	Network: "",
	Addresses: map[string]bool{
		ens:  true,
		tusd: true,
	},
	Abis: map[string]string{
		ens:  "",
		tusd: "",
	},
	Events: map[string][]string{
		ens:  {"NewOwner"},
		tusd: {"Transfer"},
	},
	Methods: map[string][]string{
		ens:  nil,
		tusd: nil,
	},
	MethodArgs: map[string][]string{
		ens:  nil,
		tusd: nil,
	},
	EventArgs: map[string][]string{
		ens:  nil,
		tusd: nil,
	},
	StartingBlocks: map[string]int64{
		ens:  3327417,
		tusd: 5197514,
	},
}

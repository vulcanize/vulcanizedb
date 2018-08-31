// Copyright 2018 Vulcanize
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

package generic

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/vulcanize/vulcanizedb/examples/constants"
)

type ContractConfig struct {
	Address    string
	Owner      string
	Abi        string
	ParsedAbi  abi.ABI
	FirstBlock int64
	LastBlock  int64
	Name       string
}

var DaiConfig = ContractConfig{
	Address:    constants.DaiContractAddress,
	Owner:      constants.DaiContractOwner,
	Abi:        constants.DaiAbiString,
	FirstBlock: int64(4752008),
	LastBlock:  -1,
	Name:       "Dai",
}

var TusdConfig = ContractConfig{
	Address:    constants.TusdContractAddress,
	Owner:      constants.TusdContractOwner,
	Abi:        constants.TusdAbiString,
	FirstBlock: int64(5197514),
	LastBlock:  -1,
	Name:       "Tusd",
}

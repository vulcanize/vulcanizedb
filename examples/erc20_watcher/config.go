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

package erc20_watcher

import "github.com/vulcanize/vulcanizedb/examples/constants"

type ContractConfig struct {
	Address    string
	Abi        string
	FirstBlock int64
	LastBlock  int64
	Name       string
}

var DaiConfig = ContractConfig{
	Address:    constants.DaiContractAddress,
	Abi:        constants.DaiAbiString,
	FirstBlock: int64(4752008),
	LastBlock:  -1,
	Name:       "Dai",
}

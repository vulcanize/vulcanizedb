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

package flip_kick

type TransformerConfig struct {
	ContractAddress     string
	ContractAbi         string
	Topics              []string
	StartingBlockNumber int64
	EndingBlockNumber   int64
}

var FlipKickConfig = TransformerConfig{
	ContractAddress:     "0x08cb6176addcca2e1d1ffe21bee464b72ee4cd8d", //this is a temporary address deployed locally
	ContractAbi:         FlipperABI,
	Topics:              []string{FlipKickSignature},
	StartingBlockNumber: 0,
	EndingBlockNumber:   100,
}

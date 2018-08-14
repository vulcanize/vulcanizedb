// Copyright © 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package price_feeds

var (
	PepAddress = "0x99041F808D598B782D5a3e498681C2452A31da08"
	PipAddress = "0x729D19f657BD0614b4985Cf1D82531c67569197B"
	RepAddress = "0xF5f94b7F9De14D43112e713835BCef2d55b76c1C"
)

type IPriceFeedConfig struct {
	ContractAddresses   []string
	StartingBlockNumber int64
	EndingBlockNumber   int64
}

var PriceFeedConfig = IPriceFeedConfig{
	ContractAddresses: []string{
		PepAddress,
		PipAddress,
		RepAddress,
	},
	StartingBlockNumber: 0,
	EndingBlockNumber:   100,
}

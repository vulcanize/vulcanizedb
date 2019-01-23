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

package tend

import (
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
)

func GetTendConfig() shared.TransformerConfig {
	return shared.TransformerConfig{
		TransformerName:     constants.TendLabel,
		ContractAddresses:   []string{constants.FlapperContractAddress(), constants.FlipperContractAddress()},
		ContractAbi:         constants.FlipperABI(),
		Topic:               constants.GetTendFunctionSignature(),
		StartingBlockNumber: shared.MinInt64([]int64{constants.FlapperDeploymentBlock(), constants.FlipperDeploymentBlock()}),
		EndingBlockNumber:   -1,
	}
}

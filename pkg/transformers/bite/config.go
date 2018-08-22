/*
 *  Copyright 2018 Vulcanize
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package bite

import (
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
)

var BiteConfig = shared.TransformerConfig{
	ContractAddresses:   "0xe0f0fa6982c59d8aa4ae0134bfe048327bd788cacf758b643ca41f055ffce76c", //this is a temporary address deployed locally
	ContractAbi:         BiteABI,
	Topics:              []string{BiteSignature},
	StartingBlockNumber: 0,
	EndingBlockNumber:   100,
}

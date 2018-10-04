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

package shared

import "math/big"

func BigIntToInt64(value *big.Int) int64 {
	if value == nil {
		return int64(0)
	} else {
		return value.Int64()
	}
}

func BigIntToString(value *big.Int) string {
	result := value.String()
	if result == "<nil>" {
		return ""
	} else {
		return result
	}
}

func GetDataBytesAtIndex(n int, logData []byte) []byte {
	switch {
	case n == -1:
		return logData[len(logData)-DataItemLength:]
	case n == -2:
		return logData[len(logData)-(2*DataItemLength) : len(logData)-DataItemLength]
	case n == -3:
		return logData[len(logData)-(3*DataItemLength) : len(logData)-(2*DataItemLength)]
	}
	return []byte{}
}

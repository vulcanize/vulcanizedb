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

import (
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
	"math/big"
)

var (
	rayBase      = big.NewFloat(1e27)
	wadBase      = big.NewFloat(1e18)
	rayPrecision = 27
	wadPrecision = 18
	ray          = "ray"
	wad          = "wad"
)

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
		return logData[len(logData)-constants.DataItemLength:]
	case n == -2:
		return logData[len(logData)-(2*constants.DataItemLength) : len(logData)-constants.DataItemLength]
	case n == -3:
		return logData[len(logData)-(3*constants.DataItemLength) : len(logData)-(2*constants.DataItemLength)]
	}
	return []byte{}
}

func ConvertToRay(value string) string {
	return convert(ray, value, rayPrecision)
}

func ConvertToWad(value string) string {
	return convert(wad, value, wadPrecision)
}

func convert(conversion string, value string, precision int) string {
	result := big.NewFloat(0.0)
	bigFloat := big.NewFloat(0.0)
	bigFloat.SetString(value)

	switch conversion {
	case ray:
		result.Quo(bigFloat, rayBase)
	case wad:
		result.Quo(bigFloat, wadBase)
	}
	return result.Text('f', precision)
}

func MinInt64(ints []int64) (min int64) {
	if len(ints) == 0 {
		return 0
	}
	min = ints[0]
	for _, i := range ints {
		if i < min {
			min = i
		}
	}
	return
}

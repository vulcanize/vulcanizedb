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

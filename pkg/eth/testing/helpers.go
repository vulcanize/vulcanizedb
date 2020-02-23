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

package testing

import (
	"github.com/sirupsen/logrus"
	"github.com/vulcanize/vulcanizedb/pkg/eth"
	"github.com/vulcanize/vulcanizedb/pkg/eth/core"
	"github.com/vulcanize/vulcanizedb/test_config"
)

func SampleContract() core.Contract {
	return core.Contract{
		Abi:  sampleAbiFileContents(),
		Hash: "0xd26114cd6EE289AccF82350c8d8487fedB8A0C07",
	}
}

func sampleAbiFileContents() string {
	abiFileContents, err := eth.ReadAbiFile(test_config.ABIFilePath + "sample_abi.json")
	if err != nil {
		logrus.Fatal(err)
	}
	return abiFileContents
}

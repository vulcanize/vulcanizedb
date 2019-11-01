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

package constants

import (
	"fmt"
	"math"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var initialized = false

func initConfig() {
	if initialized {
		return
	}

	if err := viper.ReadInConfig(); err == nil {
		log.Info("Using config file:", viper.ConfigFileUsed())
	} else {
		panic(fmt.Sprintf("Could not find environment file: %v", err))
	}
	initialized = true
}

// GetMinDeploymentBlock gets the minimum deployment block for multiple contracts from config
func GetMinDeploymentBlock() uint64 {
	initConfig()
	contractNames := getContractNames()
	if len(contractNames) < 1 {
		log.Fatalf("No contracts supplied")
	}
	minBlock := uint64(math.MaxUint64)
	for _, c := range contractNames {
		deployed := getDeploymentBlock(c)
		if deployed < minBlock {
			minBlock = deployed
		}
	}
	return minBlock
}

func getContractNames() []string {
	initConfig()
	transformerNames := viper.GetStringSlice("exporter.transformerNames")
	contractNames := make([]string, 0)
	for _, transformerName := range transformerNames {
		configKey := "exporter." + transformerName + ".contracts"
		names := viper.GetStringSlice(configKey)
		for _, name := range names {
			contractNames = appendNoDuplicates(transformerNames, name)
		}
	}
	return contractNames
}

func appendNoDuplicates(strSlice []string, str string) []string {
	for _, strInSlice := range strSlice {
		if strInSlice == str {
			return strSlice
		}
	}
	return append(strSlice, str)
}

func getDeploymentBlock(contractName string) uint64 {
	configKey := "contract." + contractName + ".deployed"
	value := viper.GetInt64(configKey)
	if value < 0 {
		log.Infof("No deployment block configured for contract \"%v\", defaulting to 0.", contractName)
		return 0
	}
	return uint64(value)
}

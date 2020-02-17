// Copyright © 2020 NAME HERE <EMAIL ADDRESS>
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

package cmd

import (
	"fmt"

	"github.com/makerdao/vulcanizedb/libraries/shared/storage/backfill"
	"github.com/makerdao/vulcanizedb/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	getStorageValueBlockNumber int64
	getStorageFlagName         = "get-storage-value-block-number"
)

// getStorageValueCmd represents the getStorageValue command
var getStorageValueCmd = &cobra.Command{
	Use:   "getStorageValue",
	Short: "Gets all storage values for configured contracts at the given block.",
	Long: fmt.Sprintf(`Fetches and persists storage values of the configured contracts at a given block. It is important to note that the storage value gotten with this command may not be different from the previous block in the database.

Use: ./vulcanizedb getStorageValueAt --config=<config.toml> --%s=<block number>`, getStorageFlagName),
	RunE: func(cmd *cobra.Command, args []string) error {
		SubCommand = cmd.CalledAs()
		LogWithCommand = *logrus.WithField("SubCommand", SubCommand)
		LogWithCommand.Infof("Getting storage values for all known keys at block %d", getStorageValueBlockNumber)

		validationErr := validateBlockNumberArg(getStorageValueBlockNumber, getStorageFlagName)
		if validationErr != nil {
			return validationErr
		}

		getStorageErr := getStorageAt(getStorageValueBlockNumber)
		if getStorageErr != nil {
			return fmt.Errorf("SubCommand %v: Failed to get storage values at block %v. Err: %v", SubCommand, getStorageValueBlockNumber, getStorageErr)
		}
		return nil
	},
}

func init() {
	getStorageValueCmd.Flags().Int64VarP(&getStorageValueBlockNumber, getStorageFlagName, "b", -1, "block number to fetch storage at for all configured transformers")
	rootCmd.AddCommand(getStorageValueCmd)
}

func getStorageAt(blockNumber int64) error {
	blockChain := getBlockChain()
	db := utils.LoadPostgres(databaseConfig, blockChain.Node())
	_, storageInitializers, _, exportTransformersErr := exportTransformers()
	if exportTransformersErr != nil {
		return fmt.Errorf("SubCommand %v: exporting transformers failed: %v", SubCommand, exportTransformersErr)
	}

	if len(storageInitializers) == 0 {
		return fmt.Errorf("SubCommand %v: no storage transformers found in the given config", SubCommand)
	}

	loader := backfill.NewStorageValueLoader(blockChain, &db, storageInitializers, blockNumber)
	return loader.Run()
}

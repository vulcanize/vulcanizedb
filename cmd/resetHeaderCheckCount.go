// VulcanizeDB
// Copyright © 2020 elizabethengelman

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

package cmd

import (
	"errors"
	"fmt"

	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var blockNumber int64

// resetHeaderCheckCountCmd represents the resetHeaderCheckCount command
var resetHeaderCheckCountCmd = &cobra.Command{
	Use:   "resetHeaderCheckCount",
	Short: "Resets header check_count for the given block number",
	Long: `Resets check_count to zero for the given header so that the execute command may recheck that header's logs in case one was missed.

Use: ./vulcanizedb resetHeaderCheckCount --header-block-number=<block number>
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		SubCommand = cmd.CalledAs()
		LogWithCommand = *logrus.WithField("SubCommand", SubCommand)
		LogWithCommand.Infof("Updating check_count for header %v set to 0.", blockNumber)

		validationErr := validateBlockNumberArg(blockNumber)
		if validationErr != nil {
			errorString := fmt.Sprintf("%v: header-block-number argument is required and not value was given.", SubCommand)
			return errors.New(errorString)
		}

		resetErr := resetHeaderCount(int64(blockNumber))
		if resetErr != nil {
			errorString := fmt.Sprintf("%v: Failed to reset header %v check_count to 0. Err: %v", SubCommand, blockNumber, resetErr)
			return errors.New(errorString)
		}

		return nil
	},
}

func init() {
	resetHeaderCheckCountCmd.Flags().Int64VarP(&blockNumber, "header-block-number", "b", -1, "block number of the header check_count to reset")
	rootCmd.AddCommand(resetHeaderCheckCountCmd)
}

func validateBlockNumberArg(blockNumber int64) error {
	if blockNumber == -1 {
		return errors.New("header-block-number argument is required and no value was given")
	}
	return nil
}

func resetHeaderCount(blockNumber int64) error {
	blockChain := getBlockChain()
	db := utils.LoadPostgres(databaseConfig, blockChain.Node())
	repo := repositories.NewCheckedHeadersRepository(&db)
	return repo.MarkSingleHeaderUnchecked(blockNumber)
}

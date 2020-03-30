package cmd

import (
	"fmt"

	"github.com/makerdao/vulcanizedb/libraries/shared/storage/backfill"
	"github.com/makerdao/vulcanizedb/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	backfillStorageStartBlockNumber int64
	backfillStorageEndBlockNumber   int64
	backfillStorageStartBlockFlag   = "backfill-storage-start-block"
	backfillStorageEndBlockFlag     = "backfill-storage-end-block"
)

// backfillStorageCmd represents the backfillStorage command
var backfillStorageCmd = &cobra.Command{
	Use:   "backfillStorage",
	Short: "Backfill smart contract storage for a range of blocks",
	Long: `Fetch and persist contract storage for a range of blocks. Useful if you have started watching diffs from a
new contract but do not have storage data from before you started running the transformer.
Requires a config file structured the same as it would be for running compose or composeAndExecute (to specify which
addresses and storage slots must be back-filled).
Requires values for backfill-storage-start-block (-s) and backfill-storage-end-block (-e) to define range of blocks that
need to be back-filled.
Before running this command, verify that you have run headerSync and execute for the desired blocks. Headers are
required for generating queries for storage slots by hash, and execute is required since the identifier for storage
slots that represent mappings and dynamic arrays depend on data derived from events.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		SubCommand = cmd.CalledAs()
		LogWithCommand = *logrus.WithField("SubCommand", SubCommand)
		return backfillStorage()
	},
}

func init() {
	rootCmd.AddCommand(backfillStorageCmd)
	backfillStorageCmd.Flags().Int64VarP(&backfillStorageStartBlockNumber, backfillStorageStartBlockFlag, "s", -1, "starting block from which to back-fill storage")
	backfillStorageCmd.Flags().Int64VarP(&backfillStorageEndBlockNumber, backfillStorageEndBlockFlag, "e", -1, "ending block for back-filling storage")
}

func backfillStorage() error {
	validationErr := validateBackfillStorageArgs()
	if validationErr != nil {
		return validationErr
	}

	blockChain := getBlockChain()
	db := utils.LoadPostgres(databaseConfig, blockChain.Node())

	_, storageInitializers, _, exportTransformersErr := exportTransformers()
	if exportTransformersErr != nil {
		return fmt.Errorf("SubCommand %v: exporting transformers failed: %v", SubCommand, exportTransformersErr)
	}

	if len(storageInitializers) == 0 {
		return fmt.Errorf("SubCommand %v: no storage transformers found in the given config", SubCommand)
	}

	loader := backfill.NewStorageValueLoader(blockChain, &db, storageInitializers, backfillStorageStartBlockNumber, backfillStorageEndBlockNumber)
	return loader.Run()
}

func validateBackfillStorageArgs() error {
	validateStartBlockErr := validateBlockNumberArg(backfillStorageStartBlockNumber, backfillStorageStartBlockFlag)
	if validateStartBlockErr != nil {
		return validateStartBlockErr
	}

	validateEndBlockErr := validateBlockNumberArg(backfillStorageEndBlockNumber, backfillStorageEndBlockFlag)
	if validateEndBlockErr != nil {
		return validateEndBlockErr
	}

	return nil
}
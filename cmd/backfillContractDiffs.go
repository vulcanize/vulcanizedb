package cmd

import (
	"fmt"

	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/backfill"
	"github.com/makerdao/vulcanizedb/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	backfillContractDiffsAddress     string
	backfillContractDiffsBlockNumber int64
	backfillContractDiffsAddressFlag = "all-contract-values-address"
	backfillContractDiffsBlockFlag   = "all-contract-values-block-number"
)

var backfillContractDiffsCmd = &cobra.Command{
	Use:   "backfillContractDiffs",
	Short: "Backfills all diffs for the given contract starting at the given block.",
	Long: fmt.Sprintf(`Backfills diffs belonging to the given contract starting at the given block.

Use: ./vulcanizedb backfillContractDiffs --config=<config.toml> --%s=<block number> --%s=<address>`, backfillContractDiffsBlockFlag, backfillContractDiffsAddressFlag),
	RunE: func(cmd *cobra.Command, args []string) error {
		SubCommand = cmd.CalledAs()
		LogWithCommand = *logrus.WithField("SubCommand", SubCommand)
		LogWithCommand.Infof("Backfilling diffs for contract %d", backfillContractDiffsBlockNumber)

		if backfillContractDiffsAddress == "" {
			return fmt.Errorf("SubCommand %v: %s argument is required and no value was given", SubCommand, backfillContractDiffsAddressFlag)
		}

		validationErr := validateBlockNumberArg(backfillContractDiffsBlockNumber, backfillContractDiffsBlockFlag)
		if validationErr != nil {
			return validationErr
		}

		backfillDiffsErr := backfillContractDiffs(backfillContractDiffsAddress, backfillContractDiffsBlockNumber)
		if backfillDiffsErr != nil {
			return fmt.Errorf("SubCommand %v: Failed to backfill diffs for contract %v. Err: %v", SubCommand, backfillContractDiffsAddress, backfillDiffsErr)
		}
		return nil
	},
}

func init() {
	backfillContractDiffsCmd.Flags().StringP(backfillContractDiffsAddressFlag, "a", "", "address whose diffs will be backfilled")
	backfillContractDiffsCmd.Flags().Int64VarP(&backfillContractDiffsBlockNumber, backfillContractDiffsBlockFlag, "b", -1, "block number at which to start backfilling diffs for the given contract")
	rootCmd.AddCommand(backfillContractDiffsCmd)
}

func backfillContractDiffs(address string, startingBlock int64) error {
	blockChain := getBlockChain()
	db := utils.LoadPostgres(databaseConfig, blockChain.Node())
	_, storageInitializers, _, exportTransformersErr := exportTransformers()
	if exportTransformersErr != nil {
		return fmt.Errorf("SubCommand %v: exporting transformers failed: %v", SubCommand, exportTransformersErr)
	}

	if len(storageInitializers) == 0 {
		return fmt.Errorf("SubCommand %v: no storage transformers found in the given config", SubCommand)
	}

	latestBlock, latestBlockErr := blockChain.LastBlock()
	if latestBlockErr != nil {
		return fmt.Errorf("SubCommand %v: error getting latest block: %v", SubCommand, latestBlockErr)
	}

	filteredInitializers, filterErr := filterByAddress(address, storageInitializers)
	if filterErr != nil {
		return filterErr
	}

	loader := backfill.NewStorageValueLoader(blockChain, &db, filteredInitializers, startingBlock, latestBlock.Int64())
	return loader.Run()
}

func filterByAddress(address string, initializers []storage.TransformerInitializer) ([]storage.TransformerInitializer, error) {
	for _, initializer := range initializers {
		if initializer(nil).GetContractAddress().Hex() == address {
			return []storage.TransformerInitializer{initializer}, nil
		}
	}
	return nil, fmt.Errorf("subcommand %v: no storage transformer found with address %v", SubCommand, address)
}

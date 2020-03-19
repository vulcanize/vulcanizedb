package cmd

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/backfill"
	"github.com/makerdao/vulcanizedb/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	backfillContractDiffsAddress          string
	backfillContractDiffsStartBlockNumber int64
	backfillContractDiffsEndBlockNumber   int64
	backfillContractDiffsAddressFlag      = "all-contract-values-address"
	backfillContractDiffsStartBlockFlag   = "all-contract-values-start-block"
	backfillContractDiffsEndBlockFlag     = "all-contract-values-end-block"
)

var backfillContractDiffsCmd = &cobra.Command{
	Use:   "backfillContractDiffs",
	Short: "Backfills all diffs for the given contract, starting at the given start-block, and optionally ending at the given end-block.",
	Long: fmt.Sprintf(`Backfills all diffs for the given contract, starting at the given start-block, and optionally ending at the given end-block.

Use: ./vulcanizedb backfillContractDiffs --config=<config.toml> --%s=<address> --%s=<block number> [ --%s=<block number> ]`, backfillContractDiffsAddressFlag, backfillContractDiffsStartBlockFlag, backfillContractDiffsEndBlockFlag),
	RunE: func(cmd *cobra.Command, args []string) error {
		SubCommand = cmd.CalledAs()
		LogWithCommand = *logrus.WithField("SubCommand", SubCommand)
		LogWithCommand.Infof("Backfilling diffs for contract %v", backfillContractDiffsAddress)

		if backfillContractDiffsAddress == "" {
			return fmt.Errorf("SubCommand %v: %s argument is required and no value was given", SubCommand, backfillContractDiffsAddressFlag)
		}

		validationErr := validateBlockNumberArg(backfillContractDiffsStartBlockNumber, backfillContractDiffsStartBlockFlag)
		if validationErr != nil {
			return validationErr
		}

		backfillDiffsErr := backfillContractDiffs(backfillContractDiffsAddress, backfillContractDiffsStartBlockNumber, backfillContractDiffsEndBlockNumber)
		if backfillDiffsErr != nil {
			return fmt.Errorf("SubCommand %v: Failed to backfill diffs for contract %v. Err: %v", SubCommand, backfillContractDiffsAddress, backfillDiffsErr)
		}
		return nil
	},
}

func init() {
	backfillContractDiffsCmd.Flags().StringVarP(&backfillContractDiffsAddress, backfillContractDiffsAddressFlag, "a", "", "address whose diffs will be backfilled")
	backfillContractDiffsCmd.Flags().Int64VarP(&backfillContractDiffsStartBlockNumber, backfillContractDiffsStartBlockFlag, "s", -1, "block number at which to start backfilling diffs for the given contract")
	backfillContractDiffsCmd.Flags().Int64VarP(&backfillContractDiffsEndBlockNumber, backfillContractDiffsEndBlockFlag, "e", -1, "optional block number at which to stop backfilling diffs for the given contract (defaults to latest block)")
	rootCmd.AddCommand(backfillContractDiffsCmd)
}

func backfillContractDiffs(address string, startingBlock, endingBlock int64) error {
	blockChain := getBlockChain()
	db := utils.LoadPostgres(databaseConfig, blockChain.Node())
	_, storageInitializers, _, exportTransformersErr := exportTransformers()
	if exportTransformersErr != nil {
		return fmt.Errorf("SubCommand %v: exporting transformers failed: %v", SubCommand, exportTransformersErr)
	}

	if len(storageInitializers) == 0 {
		return fmt.Errorf("SubCommand %v: no storage transformers found in the given config", SubCommand)
	}

	if endingBlock == -1 {
		latestBlock, latestBlockErr := blockChain.LastBlock()
		if latestBlockErr != nil {
			return fmt.Errorf("SubCommand %v: error getting latest block: %v", SubCommand, latestBlockErr)
		}
		endingBlock = latestBlock.Int64()
	}

	filteredInitializers, filterErr := filterByAddress(address, storageInitializers)
	if filterErr != nil {
		return filterErr
	}

	loader := backfill.NewStorageValueLoader(blockChain, &db, filteredInitializers, startingBlock, endingBlock)
	return loader.Run()
}

func filterByAddress(address string, initializers []storage.TransformerInitializer) ([]storage.TransformerInitializer, error) {
	for _, initializer := range initializers {
		if initializer(nil).GetContractAddress() == common.HexToAddress(address) {
			return []storage.TransformerInitializer{initializer}, nil
		}
	}
	return nil, fmt.Errorf("subcommand %v: no storage transformer found with address %v", SubCommand, address)
}

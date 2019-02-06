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

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/fs"
	"github.com/vulcanize/vulcanizedb/pkg/transformers"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/storage"
	"log"
)

// parseStorageDiffsCmd represents the parseStorageDiffs command
var parseStorageDiffsCmd = &cobra.Command{
	Use:   "parseStorageDiffs",
	Short: "Continuously ingest storage diffs from a CSV file",
	Long: `Read storage diffs out of a CSV file that is constantly receiving
new rows from an Ethereum node. For example:

./vulcanizedb parseStorageDiffs --config=environments/staging.toml

Note that the path to your storage diffs must be configured in your toml
file under storageDiffsPath.`,
	Run: func(cmd *cobra.Command, args []string) {
		parseStorageDiffs()
	},
}

func init() {
	rootCmd.AddCommand(parseStorageDiffsCmd)
}

func parseStorageDiffs() {
	blockChain := getBlockChain()
	db, err := postgres.NewDB(databaseConfig, blockChain.Node())
	if err != nil {
		log.Fatal("Failed to initialize database: ", err)
	}

	tailer := fs.FileTailer{Path: storageDiffsPath}

	// TODO: configure transformers
	watcher := shared.NewStorageWatcher(tailer, db)
	watcher.AddTransformers([]storage.TransformerInitializer{
		transformers.GetPitStorageTransformer().NewTransformer,
		transformers.GetVatStorageTransformer().NewTransformer,
	})

	err = watcher.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

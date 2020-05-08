// Copyright © 2020 Vulcanize, Inc
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"

	"github.com/spf13/cobra"
)

// validateRootNodeCmd represents the validateRootNode command
var validateRootNodeCmd = &cobra.Command{
	Use:   "validateRootNode",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("validateRootNode called")
		validateRootNode()
	},
}

func validateRootNode() {
	dbConfig := new(config.Database)
	dbConfig.Init()
	connectString := dbConfig.ConnectionString()
	db, err := sqlx.Connect("postgres", connectString)
	if err != nil {
		logrus.Fatal(err)
	}
	tx, err := db.Beginx()
	if err != nil {
		logrus.Fatal(err)
	}
	pgStr := `SELECT state_cids.cid FROM eth.state_cids 
			INNER JOIN eth.header_cids ON (state_cids.header_id = header_cids.id)
			WHERE block_number = $1 AND state_path = $2`
	var cidString string
	if err := tx.Get(&cidString, pgStr, blockNumber, []byte{}); err != nil {
		shared.Rollback(tx)
		logrus.Fatal(err)
	}
	rootNode, err := shared.FetchIPLD(tx, cidString)
	if err != nil {
		shared.Rollback(tx)
		logrus.Fatal(err)
	}
	rootNodeHash := crypto.Keccak256Hash(rootNode)
	pgStr = `SELECT state_root FROM eth.header_cids
			WHERE block_number = $1`
	var stateRootString string
	if err := tx.Get(&stateRootString, pgStr, blockNumber); err != nil {
		shared.Rollback(tx)
		logrus.Fatal(err)
	}
	stateRoot := common.HexToHash(stateRootString)
	if !bytes.Equal(stateRoot.Bytes(), rootNodeHash.Bytes()) {
		logrus.Fatal("root node hash does not match state root found in header")
	}
	logrus.Info("root node hash matches state root found in header")
}

func init() {
	rootCmd.AddCommand(validateRootNodeCmd)
	validateRootNodeCmd.Flags().IntVarP(&blockNumber, "block-number", "b", 0, "Block number to write to disk")
}

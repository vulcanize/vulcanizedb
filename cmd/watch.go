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
	"fmt"

	"github.com/spf13/cobra"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch and transform data from a chain source",
	Long: `This command allows one to configure a set of wasm functions and SQL trigger functions
that call them to watch and transform data from the specified chain source.

A watcher is composed of four parts:
1) Go execution engine- this command- which fetches raw chain data and adds it to the Postres queued ready data tables
2) TOML config file which specifies what subset of chain data to fetch and from where and contains references to the below
3) Set of WASM binaries which are loaded into Postgres and used by
4) Set of PostgreSQL trigger functions which automatically act on data as it is inserted into the queued ready data tables`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("watch called")
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
}

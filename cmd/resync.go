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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/resync"
)

// resyncCmd represents the resync command
var resyncCmd = &cobra.Command{
	Use:   "resync",
	Short: "Resync historical data",
	Long:  `Use this command to fill in sections of missing data in the super node`,
	Run: func(cmd *cobra.Command, args []string) {
		subCommand = cmd.CalledAs()
		logWithCommand = *log.WithField("SubCommand", subCommand)
		rsyncCmdCommand()
	},
}

func init() {
	rootCmd.AddCommand(resyncCmd)
}

func rsyncCmdCommand() {
	rConfig, err := resync.NewReSyncConfig()
	if err != nil {
		logWithCommand.Fatal(err)
	}
	if err := ipfs.InitIPFSPlugins(); err != nil {
		logWithCommand.Fatal(err)
	}
	rService, err := resync.NewResyncService(rConfig)
	if err != nil {
		logWithCommand.Fatal(err)
	}
	if err := rService.Resync(); err != nil {
		logWithCommand.Fatal(err)
	}
	logWithCommand.Infof("%s %s resync finished", rConfig.Chain.String(), rConfig.ResyncType.String())
}

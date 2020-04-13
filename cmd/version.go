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

	v "github.com/vulcanize/vulcanizedb/version"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version of vulcanizeDB",
	Long: `Use this command to fetch the version of vulcanizeDB

Usage: ./vulcanizedb version`,
	Run: func(cmd *cobra.Command, args []string) {
		subCommand = cmd.CalledAs()
		logWithCommand = *log.WithField("SubCommand", subCommand)
		logWithCommand.Infof("VulcanizeDB version: %s", v.VersionWithMeta)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

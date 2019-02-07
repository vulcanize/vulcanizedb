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

package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/vulcanize/vulcanizedb/pkg/plugin/helpers"
)

type Plugin struct {
	Initializers map[string]string     // Map of import aliases to transformer initializer paths
	Dependencies map[string]string     // Map of vendor dep names to their repositories
	Migrations   map[string]string     // Map of vendor dep names to relative path from repository to db migrations
	Types        map[string]PluginType // Map of import aliases to their transformer initializer type (e.g. eth-event vs eth-storage)
	FilePath     string
	FileName     string
	Save         bool
}

func (c *Plugin) GetPluginPaths() (string, string, error) {
	path, err := helpers.CleanPath(c.FilePath)
	if err != nil {
		return "", "", err
	}

	name := strings.Split(c.FileName, ".")[0]
	goFile := filepath.Join(path, name+".go")
	soFile := filepath.Join(path, name+".so")

	return goFile, soFile, nil
}

func (c *Plugin) GetMigrationsPaths() ([]string, error) {
	paths := make([]string, 0, len(c.Migrations))
	for key, relPath := range c.Migrations {
		repo, ok := c.Dependencies[key]
		if !ok {
			return nil, errors.New(fmt.Sprintf("migration %s with path %s missing repository", key, relPath))
		}
		path := filepath.Join("$GOPATH/src/github.com/vulcanize/vulcanizedb/vendor", repo, relPath)
		cleanPath, err := helpers.CleanPath(path)
		if err != nil {
			return nil, err
		}
		paths = append(paths, cleanPath)
	}

	return paths, nil
}

type PluginType int

const (
	UnknownTransformerType PluginType = iota + 1
	EthEvent
	EthStorage
	IpfsEvent
	IpfsStorage
)

func (pt PluginType) String() string {
	names := [...]string{
		"eth_event",
		"eth_storage",
		"ipfs_event",
		"ipfs_storage",
	}

	if pt > IpfsStorage || pt < EthEvent {
		return "Unknown"
	}

	return names[pt]
}

func GetPluginType(str string) PluginType {
	types := [...]PluginType{
		EthEvent,
		EthStorage,
		IpfsEvent,
		IpfsStorage,
	}

	for _, ty := range types {
		if ty.String() == str && ty.String() != "Unknown" {
			return ty
		}
	}

	return UnknownTransformerType
}

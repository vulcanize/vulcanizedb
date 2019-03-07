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
	Transformers map[string]Transformer
	FilePath     string
	FileName     string
	Save         bool
	Home         string
}

type Transformer struct {
	Path           string
	Type           TransformerType
	MigrationPath  string
	MigrationRank  int
	RepositoryPath string
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

// Removes duplicate migration paths and returns them in ranked order
func (c *Plugin) GetMigrationsPaths() ([]string, error) {
	paths := make(map[int]string)
	for name, transformer := range c.Transformers {
		repo := transformer.RepositoryPath
		mig := transformer.MigrationPath
		path := filepath.Join("$GOPATH/src", c.Home, "vendor", repo, mig)
		cleanPath, err := helpers.CleanPath(path)
		if err != nil {
			return nil, err
		}
		// If there is a different path with the same rank then we have a conflict
		_, ok := paths[transformer.MigrationRank]
		if ok {
			conflictingPath := paths[transformer.MigrationRank]
			if conflictingPath != cleanPath {
				return nil, errors.New(fmt.Sprintf("transformer %s has the same migration rank (%d) as another transformer", name, transformer.MigrationRank))
			}
		}
		paths[transformer.MigrationRank] = cleanPath
	}

	sortedPaths := make([]string, len(paths))
	for rank, path := range paths {
		sortedPaths[rank] = path
	}

	return sortedPaths, nil
}

// Removes duplicate repo paths before returning them
func (c *Plugin) GetRepoPaths() map[string]bool {
	paths := make(map[string]bool)
	for _, transformer := range c.Transformers {
		paths[transformer.RepositoryPath] = true
	}

	return paths
}

type TransformerType int

const (
	UnknownTransformerType TransformerType = iota
	EthEvent
	EthStorage
)

func (pt TransformerType) String() string {
	names := [...]string{
		"Unknown",
		"eth_event",
		"eth_storage",
	}

	if pt > EthStorage || pt < EthEvent {
		return "Unknown"
	}

	return names[pt]
}

func GetTransformerType(str string) TransformerType {
	types := [...]TransformerType{
		EthEvent,
		EthStorage,
	}

	for _, ty := range types {
		if ty.String() == str {
			return ty
		}
	}

	return UnknownTransformerType
}

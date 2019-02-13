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
	"path/filepath"
	"strings"

	"github.com/vulcanize/vulcanizedb/pkg/plugin/helpers"
)

type Plugin struct {
	Transformers map[string]Transformer
	FilePath     string
	FileName     string
	Save         bool
}

type Transformer struct {
	Path           string
	Type           TransformerType
	MigrationPath  string
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

// Removes duplicate migration paths before returning them
func (c *Plugin) GetMigrationsPaths() (map[string]bool, error) {
	paths := make(map[string]bool)
	for _, transformer := range c.Transformers {
		repo := transformer.RepositoryPath
		mig := transformer.MigrationPath
		path := filepath.Join("$GOPATH/src/github.com/vulcanize/vulcanizedb/vendor", repo, mig)
		cleanPath, err := helpers.CleanPath(path)
		if err != nil {
			return nil, err
		}
		paths[cleanPath] = true
	}

	return paths, nil
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

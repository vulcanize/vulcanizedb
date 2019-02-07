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

package writer

import (
	"errors"
	"fmt"

	. "github.com/dave/jennifer/jen"

	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/plugin/helpers"
)

type PluginWriter interface {
	WritePlugin() error
}

type writer struct {
	GenConfig config.Plugin
}

func NewPluginWriter(gc config.Plugin) *writer {
	return &writer{
		GenConfig: gc,
	}
}

// Generates the plugin code
func (w *writer) WritePlugin() error {
	// Setup plugin file paths
	goFile, err := w.setupFilePath()
	if err != nil {
		return err
	}

	// Begin code generation
	f := NewFile("main")
	f.HeaderComment("This is a plugin generated to export the configured transformer initializers")

	// Import TransformerInitializers specified in config
	f.ImportAlias("github.com/vulcanize/vulcanizedb/libraries/shared/transformer", "interface")
	for alias, imp := range w.GenConfig.Initializers {
		f.ImportAlias(imp, alias)
	}

	// Collect initializer code
	ethEventInitializers, ethStorageInitializers, _, _ := w.sortTransformers()

	// Create Exporter variable with method to export the set of the imported storage and event transformer initializers
	f.Type().Id("exporter").String()
	f.Var().Id("Exporter").Id("exporter")
	f.Func().Params(Id("e").Id("exporter")).Id("Export").Params().Index().Qual(
		"github.com/vulcanize/vulcanizedb/libraries/shared/transformer",
		"TransformerInitializer").Index().Qual(
		"github.com/vulcanize/vulcanizedb/libraries/shared/transformer",
		"StorageTransformerInitializer").Block(
		Return(Index().Qual(
			"github.com/vulcanize/vulcanizedb/libraries/shared/transformer",
			"TransformerInitializer").Values(ethEventInitializers...)),
		Index().Qual(
			"github.com/vulcanize/vulcanizedb/libraries/shared/transformer",
			"StorageTransformerInitializer").Values(ethStorageInitializers...)) // Exports the collected initializers

	// Write code to destination file
	err = f.Save(goFile)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to save generated .go file: %s\r\n%s", goFile, err.Error()))
	}
	return nil
}

func (w *writer) sortTransformers() ([]Code, []Code, []Code, []Code) {
	// Collect code for various initializers
	importedEthEventInitializers := make([]Code, 0)
	importerEthStorageInitializers := make([]Code, 0)
	importedIpfsEventInitializers := make([]Code, 0)
	importerIpfsStorageInitializers := make([]Code, 0)
	for name, path := range w.GenConfig.Initializers {
		switch w.GenConfig.Types[name] {
		case config.EthEvent:
			importedEthEventInitializers = append(importedEthEventInitializers, Qual(path, "TransformerInitializer"))
		case config.EthStorage:
			importerEthStorageInitializers = append(importerEthStorageInitializers, Qual(path, "StorageTransformerInitializer"))
		case config.IpfsEvent:
			//importedIpfsEventInitializers = append(importedIpfsEventInitializers, Qual(path, "IpfsEventTransformerInitializer"))
		case config.IpfsStorage:
			//importerIpfsStorageInitializers = append(importerIpfsStorageInitializers, Qual(path, "IpfsStorageTransformerInitializer"))
		}
	}

	return importedEthEventInitializers,
		importerEthStorageInitializers,
		importedIpfsEventInitializers,
		importerIpfsStorageInitializers
}

func (w *writer) setupFilePath() (string, error) {
	goFile, soFile, err := w.GenConfig.GetPluginPaths()
	if err != nil {
		return "", err
	}
	// Clear .go and .so files of the same name if they exist
	return goFile, helpers.ClearFiles(goFile, soFile)
}

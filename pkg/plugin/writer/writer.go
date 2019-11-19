// VulcanizeDB
// Copyright © 2019 Vulcanize

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

	"github.com/makerdao/vulcanizedb/pkg/config"
	"github.com/makerdao/vulcanizedb/pkg/plugin/helpers"
)

// Interface for writing a .go file for a simple
// plugin that exports the set of transformer
// initializers specified in the config
type PluginWriter interface {
	WritePlugin() error
}

type writer struct {
	GenConfig config.Plugin
}

// Requires populated plugin config
func NewPluginWriter(gc config.Plugin) *writer {
	return &writer{
		GenConfig: gc,
	}
}

// Generates the plugin code according to config specification
func (w *writer) WritePlugin() error {
	// Setup plugin file paths
	goFile, err := w.setupFilePath()
	if err != nil {
		return err
	}

	// Begin code generation
	f := NewFile("main")
	f.HeaderComment("This is a plugin generated to export the configured transformer initializers")

	// Import pkgs for generic TransformerInitializer interface and specific TransformerInitializers specified in config
	f.ImportAlias("github.com/makerdao/vulcanizedb/libraries/shared/transformer", "interface")
	for name, transformer := range w.GenConfig.Transformers {
		f.ImportAlias(transformer.RepositoryPath+"/"+transformer.Path, name)
	}

	// Collect initializer code
	code, err := w.collectTransformers()
	if err != nil {
		return err
	}

	// Create Exporter variable with method to export the set of the imported storage and event transformer initializers
	f.Type().Id("exporter").String()
	f.Var().Id("Exporter").Id("exporter")
	f.Func().Params(Id("e").Id("exporter")).Id("Export").Params().Parens(List(
		Index().Qual("github.com/makerdao/vulcanizedb/libraries/shared/transformer", "EventTransformerInitializer"),
		Index().Qual("github.com/makerdao/vulcanizedb/libraries/shared/transformer", "StorageTransformerInitializer"),
		Index().Qual("github.com/makerdao/vulcanizedb/libraries/shared/transformer", "ContractTransformerInitializer"),
	)).Block(Return(
		Index().Qual(
			"github.com/makerdao/vulcanizedb/libraries/shared/transformer",
			"EventTransformerInitializer").Values(code[config.EthEvent]...),
		Index().Qual(
			"github.com/makerdao/vulcanizedb/libraries/shared/transformer",
			"StorageTransformerInitializer").Values(code[config.EthStorage]...),
		Index().Qual(
			"github.com/makerdao/vulcanizedb/libraries/shared/transformer",
			"ContractTransformerInitializer").Values(code[config.EthContract]...))) // Exports the collected event and storage transformer initializers

	// Write code to destination file
	err = f.Save(goFile)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to save generated .go file: %s\r\n%s", goFile, err.Error()))
	}
	return nil
}

// Collect code for various types of initializers
func (w *writer) collectTransformers() (map[config.TransformerType][]Code, error) {
	code := make(map[config.TransformerType][]Code)
	for _, transformer := range w.GenConfig.Transformers {
		path := transformer.RepositoryPath + "/" + transformer.Path
		switch transformer.Type {
		case config.EthEvent:
			code[config.EthEvent] = append(code[config.EthEvent], Qual(path, "EventTransformerInitializer"))
		case config.EthStorage:
			code[config.EthStorage] = append(code[config.EthStorage], Qual(path, "StorageTransformerInitializer"))
		case config.EthContract:
			code[config.EthContract] = append(code[config.EthContract], Qual(path, "ContractTransformerInitializer"))
		default:
			return nil, errors.New(fmt.Sprintf("invalid transformer type %s", transformer.Type))
		}
	}

	return code, nil
}

// Setup the .go, clear old ones if present
func (w *writer) setupFilePath() (string, error) {
	goFile, soFile, err := w.GenConfig.GetPluginPaths()
	if err != nil {
		return "", err
	}
	// Clear .go and .so files of the same name if they exist
	return goFile, helpers.ClearFiles(goFile, soFile)
}

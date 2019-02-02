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

	// Collect TransformerInitializer names
	importedInitializers := make([]Code, 0, len(w.GenConfig.Initializers))
	for _, path := range w.GenConfig.Initializers {
		importedInitializers = append(importedInitializers, Qual(path, "TransformerInitializer"))
	}

	// Create Exporter variable with method to export the set of the imported TransformerInitializers
	f.Type().Id("exporter").String()
	f.Var().Id("Exporter").Id("exporter")
	f.Func().Params(Id("e").Id("exporter")).Id("Export").Params().Index().Qual(
		"github.com/vulcanize/vulcanizedb/libraries/shared/transformer",
		"TransformerInitializer").Block(
		Return(Index().Qual(
			"github.com/vulcanize/vulcanizedb/libraries/shared/transformer",
			"TransformerInitializer").Values(importedInitializers...))) // Exports the collected TransformerInitializers

	// Write code to destination file
	err = f.Save(goFile)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to save generated .go file: %s\r\n%s", goFile, err.Error()))
	}
	return nil
}

func (w *writer) setupFilePath() (string, error) {
	goFile, soFile, err := w.GenConfig.GetPluginPaths()
	if err != nil {
		return "", err
	}
	// Clear .go and .so files of the same name if they exist
	return goFile, helpers.ClearFiles(goFile, soFile)
}

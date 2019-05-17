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

package plugin

import (
	"errors"

	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/plugin/builder"
	"github.com/vulcanize/vulcanizedb/pkg/plugin/manager"
	"github.com/vulcanize/vulcanizedb/pkg/plugin/writer"
)

// Generator is the top-level interface for creating transformer plugins
type Generator interface {
	GenerateExporterPlugin() error
}

type generator struct {
	writer.PluginWriter
	builder.PluginBuilder
	manager.MigrationManager
}

// Creates a new generator from a plugin and database config
func NewGenerator(gc config.Plugin, dbc config.Database) (*generator, error) {
	if len(gc.Transformers) < 1 {
		return nil, errors.New("plugin generator is not configured with any transformers")
	}
	return &generator{
		PluginWriter:     writer.NewPluginWriter(gc),
		PluginBuilder:    builder.NewPluginBuilder(gc),
		MigrationManager: manager.NewMigrationManager(gc, dbc),
	}, nil
}

// Generates plugin for the transformer initializers specified in the generator config
// Writes plugin code  => Sets up build environment => Builds .so file => Performs db migrations for the plugin transformers => Clean up
func (g *generator) GenerateExporterPlugin() error {
	// Use plugin writer interface to write the plugin code
	err := g.PluginWriter.WritePlugin()
	if err != nil {
		return err
	}
	// Clean up temporary files and directories when we are done
	defer g.PluginBuilder.CleanUp()
	// Use plugin builder interface to setup build environment and compile .go file into a .so file
	err = g.PluginBuilder.BuildPlugin()
	if err != nil {
		return err
	}

	// Perform db migrations for the transformers
	return g.MigrationManager.RunMigrations()
}

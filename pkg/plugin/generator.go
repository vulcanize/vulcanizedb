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

package plugin

import (
	"errors"
	"fmt"
	"os"

	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/plugin/builder"
	"github.com/vulcanize/vulcanizedb/pkg/plugin/manager"
	"github.com/vulcanize/vulcanizedb/pkg/plugin/writer"
)

type Generator interface {
	GenerateExporterPlugin() error
}

type generator struct {
	writer.PluginWriter
	builder.PluginBuilder
	manager.MigrationManager
}

func NewGenerator(gc config.Plugin, dbc config.Database) (*generator, error) {
	if len(gc.Initializers) < 1 {
		return nil, errors.New("generator needs to be configured with TransformerInitializer import paths")
	}
	if len(gc.Dependencies) < 1 {
		return nil, errors.New("generator needs to be configured with root repository path(s)")
	}
	if len(gc.Migrations) < 1 {
		fmt.Fprintf(os.Stderr, "warning: no db migration paths have been provided for the plugin transformers\r\n")
	}
	return &generator{
		PluginWriter:     writer.NewPluginWriter(gc),
		PluginBuilder:    builder.NewPluginBuilder(gc, dbc),
		MigrationManager: manager.NewMigrationManager(gc, dbc),
	}, nil
}

func (g *generator) GenerateExporterPlugin() error {
	err := g.PluginWriter.WritePlugin()
	if err != nil {
		return err
	}
	defer g.PluginBuilder.CleanUp()
	err = g.PluginBuilder.BuildPlugin()
	if err != nil {
		return err
	}

	return g.MigrationManager.RunMigrations()
}

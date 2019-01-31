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

package autogen

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	. "github.com/dave/jennifer/jen"

	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/utils"
)

type Generator interface {
	GenerateExporterPlugin() error
}

type generator struct {
	GenConfig  *Config
	DBConfig   config.Database
	tmpMigDir  string
	tmpVenDirs []string
}

func NewGenerator(gc Config, dbc config.Database) *generator {
	return &generator{
		GenConfig: &gc,
		DBConfig:  dbc,
	}
}

func (g *generator) GenerateExporterPlugin() error {
	if len(g.GenConfig.Initializers) < 1 {
		return errors.New("generator needs to be configured with TransformerInitializer import paths")
	}
	if len(g.GenConfig.Dependencies) < 1 {
		return errors.New("generator needs to be configured with root repository path(s)")
	}
	if len(g.GenConfig.Migrations) < 1 {
		fmt.Fprintf(os.Stderr, "warning: no db migration paths have been provided\r\n")
	}

	// Get plugin file paths
	goFile, soFile, err := g.GenConfig.GetPluginPaths()
	if err != nil {
		return err
	}

	// Generate Exporter code
	err = g.generateCode(goFile, soFile)
	if err != nil {
		return err
	}

	// Setup temp vendor lib and migrations directories
	err = g.setupTempDirs()
	if err != nil {
		return err
	}

	// Clear tmp files and directories when we exit
	defer g.cleanUp(goFile)

	// Build the .go file into a .so plugin
	err = exec.Command("go", "build", "-buildmode=plugin", "-o", soFile, goFile).Run()
	if err != nil {
		return errors.New(fmt.Sprintf("unable to build .so file: %s", err.Error()))
	}

	// Run migrations only after successfully building .so file
	return g.runMigrations()
}

// Generates the plugin code
func (g *generator) generateCode(goFile, soFile string) error {
	// Clear .go and .so files of the same name if they exist
	err := utils.ClearFiles(goFile, soFile)
	if err != nil {
		return err
	}
	// Begin code generation
	f := NewFile("main")
	f.HeaderComment("This exporter is generated to export the configured transformer initializers")

	// Import TransformerInitializers specified in config
	f.ImportAlias("github.com/vulcanize/vulcanizedb/libraries/shared/transformer", "interface")
	for alias, imp := range g.GenConfig.Initializers {
		f.ImportAlias(imp, alias)
	}

	// Collect TransformerInitializer names
	importedInitializers := make([]Code, 0, len(g.GenConfig.Initializers))
	for _, path := range g.GenConfig.Initializers {
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
	return f.Save(goFile)
}

// Sets up temporary vendor libs and migration directories
func (g *generator) setupTempDirs() error {
	// TODO: Less hacky way of handling plugin build deps
	dirPath, err := utils.CleanPath("$GOPATH/src/github.com/vulcanize/vulcanizedb/")
	if err != nil {
		return err
	}
	vendorPath := filepath.Join(dirPath, "vendor")

	// Keep track of where we are writing transformer vendor libs, so that we can remove them afterwards
	g.tmpVenDirs = make([]string, 0, len(g.GenConfig.Dependencies))
	// Import transformer dependencies so that we build our plugin
	for name, importPath := range g.GenConfig.Dependencies {
		index := strings.Index(importPath, "/")
		gitPath := importPath[:index] + ":" + importPath[index+1:]
		importURL := "git@" + gitPath + ".git"
		depPath := filepath.Join(vendorPath, importPath)
		err = exec.Command("git", "clone", importURL, depPath).Run()
		if err != nil {
			return errors.New(fmt.Sprintf("unable to clone %s transformer dependency: %s", name, err.Error()))
		}

		err := os.RemoveAll(filepath.Join(depPath, "vendor/"))
		if err != nil {
			return err
		}

		g.tmpVenDirs = append(g.tmpVenDirs, depPath)
	}

	// Initialize temp directory for transformer migrations
	g.tmpMigDir, err = utils.CleanPath("$GOPATH/src/github.com/vulcanize/vulcanizedb/db/plugin_migrations")
	if err != nil {
		return err
	}
	err = os.RemoveAll(g.tmpMigDir)
	if err != nil {
		return errors.New(fmt.Sprintf("unable to remove file found at %s where tmp directory needs to be written", g.tmpMigDir))
	}

	return os.Mkdir(g.tmpMigDir, os.FileMode(0777))
}

func (g *generator) runMigrations() error {
	// Get paths to db migrations
	paths, err := g.GenConfig.GetMigrationsPaths()
	if err != nil {
		return err
	}
	if len(paths) < 1 {
		return nil
	}

	// Create temporary copies of migrations to the temporary migrationDir
	// These tmps are identical except they have had `1` added in front of their unix_timestamps
	// As such, they will be ran on top of all core migrations (at least, for the next ~317 years)
	// But will still be ran in the same order relative to one another
	// TODO: Less hacky way of handing migrations
	err = g.createMigrationCopies(paths)
	if err != nil {
		return err
	}

	// Run the copied migrations
	pgStr := fmt.Sprintf("postgres://%s:%d/%s?sslmode=disable", g.DBConfig.Hostname, g.DBConfig.Port, g.DBConfig.Name)
	return exec.Command("migrate", "-path", g.tmpMigDir, "-database", pgStr, "up").Run()
}

func (g *generator) createMigrationCopies(paths []string) error {
	for _, path := range paths {
		dir, err := ioutil.ReadDir(path)
		if err != nil {
			return err
		}
		for _, file := range dir {
			if file.IsDir() || len(file.Name()) < 15 || filepath.Ext(file.Name()) != ".sql" { // (10 digit unix time stamp + x + .sql) is bare minimum
				continue
			}
			_, err := strconv.Atoi(file.Name()[:10])
			if err != nil {
				fmt.Fprintf(os.Stderr, "migration file name %s does not posses 10 digit timestamp prefix\r\n", file.Name())
				continue
			}
			src := filepath.Join(path, file.Name())
			dst := filepath.Join(g.tmpMigDir, "1"+file.Name())
			err = utils.CopyFile(src, dst)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *generator) cleanUp(goFile string) error {
	if !g.GenConfig.Save {
		err := utils.ClearFiles(goFile)
		if err != nil {
			return err
		}
	}

	for _, venDir := range g.tmpVenDirs {
		err := os.RemoveAll(venDir)
		if err != nil {
			return err
		}
	}

	return os.RemoveAll(g.tmpMigDir)
}

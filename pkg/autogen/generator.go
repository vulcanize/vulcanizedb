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
	if g.GenConfig == nil {
		return errors.New("generator needs a config file")
	}
	if g.GenConfig.FilePath == "" {
		return errors.New("generator is missing file path")
	}
	if len(g.GenConfig.Initializers) < 1 {
		return errors.New("generator needs to be configured with imports")
	}

	// Get plugin file paths
	goFile, soFile, err := g.GenConfig.GetPluginPaths()
	if err != nil {
		return err
	}

	// Clear .go and .so files of the same name if they exist (overwrite)
	err = utils.ClearFiles(goFile, soFile)
	if err != nil {
		return err
	}

	// Generate Exporter code
	err = g.generateCode(goFile)
	if err != nil {
		return err
	}

	// Setup temp vendor lib and migrations directories
	err = g.setupTempDirs()
	if err != nil {
		return err
	}
	defer g.cleanUp() // Clear these up when we are done building our plugin

	// Build the .go file into a .so plugin
	err = exec.Command("go", "build", "-buildmode=plugin", "-o", soFile, goFile).Run()
	if err != nil {
		return err
	}
	// Run migrations only after successfully building .so file
	return g.runMigrations()
}

// Generates the plugin code
func (g *generator) generateCode(goFile string) error {
	// Begin code generation
	f := NewFile("main")
	f.HeaderComment("This exporter is generated to export the configured transformer initializers")

	// Import TransformerInitializers
	f.ImportAlias("github.com/vulcanize/vulcanizedb/libraries/shared/transformer", "interface")
	for alias, imp := range g.GenConfig.Initializers {
		f.ImportAlias(imp, alias)
	}

	// Collect TransformerInitializer names
	importedInitializers := make([]Code, 0, len(g.GenConfig.Initializers))
	for _, path := range g.GenConfig.Initializers {
		importedInitializers = append(importedInitializers, Qual(path, "TransformerInitializer"))
	}

	// Create Exporter variable with method to export a set of the configured TransformerInitializers
	f.Type().Id("exporter").String()
	f.Var().Id("Exporter").Id("exporter")
	f.Func().Params(
		Id("e").Id("exporter"),
	).Id("Export").Params().Index().Qual(
		"github.com/vulcanize/vulcanizedb/libraries/shared/transformer",
		"TransformerInitializer").Block(
		Return(Index().Qual(
			"github.com/vulcanize/vulcanizedb/libraries/shared/transformer",
			"TransformerInitializer").Values(importedInitializers...))) // Exports the collected TransformerInitializers

	// Write code to destination file
	return f.Save(goFile)
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
	location := "file://" + g.tmpMigDir
	pgStr := fmt.Sprintf("postgres://%s:%d/%s?sslmode=disable up", g.DBConfig.Hostname, g.DBConfig.Port, g.DBConfig.Name)
	return exec.Command("migrate", "-source", location, pgStr).Run()
}

// Sets up temporary vendor libs and migration directories
func (g *generator) setupTempDirs() error {
	// TODO: Less hacky way of handling plugin build deps
	dirPath, err := utils.CleanPath("$GOPATH/src/github.com/vulcanize/vulcanizedb/")
	if err != nil {
		return err
	}
	vendorPath := filepath.Join(dirPath, "vendor/")

	/*
	// Keep track of where we are writing transformer vendor libs, so that we can remove them afterwards
	g.tmpVenDirs = make([]string, 0, len(g.GenConfig.Dependencies))
	// Import transformer dependencies so that we build our plugin
	for _, importPath := range g.GenConfig.Dependencies {
		importURL := "https://" + importPath + ".git"
		depPath := filepath.Join(vendorPath, importPath)
		err = exec.Command("git", "clone", importURL, depPath).Run()
		if err != nil {
			return err
		}
		err := os.RemoveAll(filepath.Join(depPath, "vendor/"))
		if err != nil {
			return err
		}
		g.tmpVenDirs = append(g.tmpVenDirs, depPath)
	}
	*/

	// Keep track of where we are writing transformer vendor libs, so that we can remove them afterwards
	g.tmpVenDirs = make([]string, 0, len(g.GenConfig.Dependencies))
	for _, importPath := range g.GenConfig.Dependencies {
		depPath := filepath.Join(vendorPath, importPath)
		g.tmpVenDirs = append(g.tmpVenDirs, depPath)
	}

	// Dep ensure to make sure vendor pkgs are in place for building the plugin
	err = exec.Command("dep", "ensure").Run()
	if err != nil {
		return errors.New("failed to vendor transformer packages required to build plugin")
	}

	// Git checkout our head-state vendor libraries
	// This is necessary because we currently need to manual edit our vendored
	// go-ethereum abi library to allow for unpacking in empty interfaces and maps
	// This can be removed once the PRs against geth merged
	err = exec.Command("git", "checkout", dirPath).Run()
	if err != nil {
		return errors.New("failed to checkout vendored go-ethereum lib")
	}

	// Initialize temp directory for transformer migrations
	g.tmpMigDir, err = utils.CleanPath("$GOPATH/src/github.com/vulcanize/vulcanizedb/db/plugin_migrations")
	if err != nil {
		return err
	}
	stat, err := os.Stat(g.tmpMigDir)
	if err == nil {
		if !stat.IsDir() {
			return errors.New(fmt.Sprintf("file %s found where directory is expected", stat.Name()))
		}
	} else if os.IsNotExist(err) {
		os.Mkdir(g.tmpMigDir, os.FileMode(0777))
	} else {
		return err
	}

	return nil
}

func (g *generator) createMigrationCopies(paths []string) error {
	for _, path := range paths {
		dir, err := ioutil.ReadDir(path)
		if err != nil {
			return err
		}
		for _, file := range dir {
			if file.IsDir() || len(file.Name()) < 15 { // (10 digit unix time stamp + x + .sql) is bare minimum
				continue
			}
			_, err := strconv.Atoi(file.Name()[:10])
			if err != nil {
				fmt.Fprintf(os.Stderr, "migration file name %s does not posses 10 digit timestamp prefix", file.Name())
				continue
			}
			if filepath.Ext(file.Name()) == "sql" {
				src := filepath.Join(path, file.Name())
				dst := filepath.Join(g.tmpMigDir, "1"+file.Name())
				err = utils.CopyFile(src, dst)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (g *generator) cleanUp() error {
	for _, venDir := range g.tmpVenDirs {
		err := os.RemoveAll(venDir)
		if err != nil {
			return err
		}
	}

	return os.RemoveAll(g.tmpMigDir)
}

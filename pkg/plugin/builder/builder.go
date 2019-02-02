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

package builder

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/plugin/helpers"
)

type PluginBuilder interface {
	BuildPlugin() error
	CleanUp() error
}

type builder struct {
	GenConfig  config.Plugin
	tmpVenDirs []string
	goFile     string
}

func NewPluginBuilder(gc config.Plugin, dbc config.Database) *builder {
	return &builder{
		GenConfig:  gc,
		tmpVenDirs: make([]string, 0, len(gc.Dependencies)),
	}
}

func (b *builder) BuildPlugin() error {
	// Get plugin .go and .so file paths
	var err error
	var soFile string
	b.goFile, soFile, err = b.GenConfig.GetPluginPaths()
	if err != nil {
		return err
	}

	// setup env to build plugin
	err = b.setupBuildEnv()
	if err != nil {
		return err
	}

	// Build the .go file into a .so plugin
	err = exec.Command("go", "build", "-buildmode=plugin", "-o", soFile, b.goFile).Run()
	if err != nil {
		return errors.New(fmt.Sprintf("unable to build .so file: %s", err.Error()))
	}
	return nil
}

// Sets up temporary vendor libs needed for plugin build
func (b *builder) setupBuildEnv() error {
	// TODO: Less hacky way of handling plugin build deps
	vendorPath, err := helpers.CleanPath("$GOPATH/src/github.com/vulcanize/vulcanizedb/vendor")
	if err != nil {
		return err
	}

	// Import transformer dependencies so that we can build our plugin
	for name, importPath := range b.GenConfig.Dependencies {
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

		b.tmpVenDirs = append(b.tmpVenDirs, depPath)
	}

	return nil
}

func (b *builder) CleanUp() error {
	if !b.GenConfig.Save {
		err := helpers.ClearFiles(b.goFile)
		if err != nil {
			return err
		}
	}

	for _, venDir := range b.tmpVenDirs {
		err := os.RemoveAll(venDir)
		if err != nil {
			return err
		}
	}

	return nil
}

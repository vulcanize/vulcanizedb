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

// Interface for compile Go code written by the
// PluginWriter into a shared object (.so file)
// which can be used loaded as a plugin
type PluginBuilder interface {
	BuildPlugin() error
	CleanUp() error
}

type builder struct {
	GenConfig    config.Plugin
	dependencies []string
	tmpVenDirs   []string // Keep track of temp vendor directories
	goFile       string   // Keep track of goFile name
}

// Requires populated plugin config
func NewPluginBuilder(gc config.Plugin) *builder {
	return &builder{
		GenConfig:  gc,
		tmpVenDirs: make([]string, 0),
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
// This is to work around a conflict between plugins and vendoring (https://github.com/golang/go/issues/20481)
func (b *builder) setupBuildEnv() error {
	// TODO: Less hacky way of handling plugin build deps
	vendorPath, err := helpers.CleanPath(filepath.Join("$GOPATH/src", b.GenConfig.Home, "vendor"))
	if err != nil {
		return err
	}

	repoPaths := b.GenConfig.GetRepoPaths()

	// Import transformer dependencies so that we can build our plugin
	for importPath := range repoPaths {
		dstPath := filepath.Join(vendorPath, importPath)
		// When testing on Travis we need to clone the libs
		if b.GenConfig.Clone {
			// And if we want to be able to work with a private repo we need  to use ssh instead of https
			// and upload a permissioned ssh key to travis before deploying tests there
			index := strings.Index(importPath, "/")
			gitPath := importPath[:index] + ":" + importPath[index+1:]
			importURL := "git@" + gitPath + ".git"
			err = exec.Command("git", "clone", importURL, dstPath).Run()
			if err != nil {
				return errors.New(fmt.Sprintf("unable to clone transformer dependency from %s to %s: %s", importPath, dstPath, err.Error()))
			}
		} else { // If not on Travis we can work with libs at $GOPATH
			srcDir, err := helpers.CleanPath(filepath.Join("$GOPATH/src", importPath))
			if err != nil {
				return err
			}
			sp := strings.Split(dstPath, "/")
			spj := strings.Join(sp[:len(sp)-1], "/")
			err = exec.Command("rsync", "-a", srcDir, spj).Run()
			if err != nil {
				return errors.New(fmt.Sprintf("unable to copy transformer dependency from %s to %s: %s", srcDir, dstPath, err.Error()))
			}
		}
		// Have to clear out the copied over vendor lib or plugin won't build (see issue above)
		err := os.RemoveAll(filepath.Join(dstPath, "vendor"))
		if err != nil {
			return err
		}
		// Keep track of this vendor directory to clear later
		b.tmpVenDirs = append(b.tmpVenDirs, dstPath)
	}

	return nil
}

// Used to clear all of the tmp vendor libs used to build the plugin
// Also clears the go file if saving it has not been specified in the config
// Do not call until after the MigrationManager has performed its operations
// as it needs to pull the db migrations from the tmpVenDirs
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

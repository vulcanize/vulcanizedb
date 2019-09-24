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

package composer

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/plugin/helpers"
)

type composer struct {
	GenConfig    config.Plugin
	dependencies []string
	tmpVenDirs   []string // Keep track of temp vendor directories
	goFile       string   // Keep track of goFile name
}

// Requires populated plugin config
func NewPluginComposer(gc config.Plugin) *composer {
	return &composer{
		GenConfig:  gc,
		tmpVenDirs: make([]string, 0),
	}
}

func (b *composer) BuildPlugin() error {
	// Get plugin .go and .so file paths
	var err error
	b.goFile, _, err = b.GenConfig.GetPluginPaths()
	if err != nil {
		return err
	}

	// setup env to build plugin
	return b.setupComposeEnv()
}

// Sets up temporary vendor libs needed for plugin build
// This is to work around a conflict between plugins and vendoring (https://github.com/golang/go/issues/20481)
func (b *composer) setupComposeEnv() error {
	// TODO: Less hacky way of handling plugin build deps
	vendorPath, err := helpers.CleanPath(filepath.Join("$GOPATH/src", b.GenConfig.Home, "vendor"))
	if err != nil {
		return err
	}

	repoPaths := b.GenConfig.GetRepoPaths()

	// Import transformer dependencies so that we can build our plugin
	for importPath := range repoPaths {
		dst := filepath.Join(vendorPath, importPath)
		src, cleanErr := helpers.CleanPath(filepath.Join("$GOPATH/src", importPath))
		if cleanErr != nil {
			return cleanErr
		}

		copyErr := helpers.CopyDir(src, dst, "vendor")
		if copyErr != nil {
			return errors.New(fmt.Sprintf("unable to copy transformer dependency from %s to %s: %v", src, dst, copyErr))
		}

		// Have to clear out the copied over vendor lib or plugin won't build (see issue above)
		removeErr := os.RemoveAll(filepath.Join(dst, "vendor"))
		if removeErr != nil {
			return removeErr
		}
		// Keep track of this vendor directory to clear later
		b.tmpVenDirs = append(b.tmpVenDirs, dst)
	}

	return nil
}

func (c *composer) CleanUp() error {
	for _, venDir := range c.tmpVenDirs {
		err := os.RemoveAll(venDir)
		if err != nil {
			return err
		}
	}
	return nil
}

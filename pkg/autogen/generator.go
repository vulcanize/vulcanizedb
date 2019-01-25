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
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	. "github.com/dave/jennifer/jen"
	"github.com/mitchellh/go-homedir"
)

type Generator interface {
	GenerateTransformerPlugin() error
}

type generator struct {
	*Config
}

func NewGenerator(config Config) *generator {
	return &generator{
		Config: &config,
	}
}

func (g *generator) GenerateTransformerPlugin() error {
	if g.Config == nil {
		return errors.New("generator needs a config file")
	}
	if g.Config.FilePath == "" {
		return errors.New("generator is missing file path")
	}
	if len(g.Config.Imports) < 1 {
		return errors.New("generator needs to be configured with imports")
	}

	// Create file path
	goFile, soFile, err := GetPaths(*g.Config)
	if err != nil {
		return err
	}

	// Clear previous .go and .so files if they exist
	err = ClearFiles(goFile, soFile)
	if err != nil {
		return err
	}

	// Begin code generation
	f := NewFile("main")
	f.HeaderComment("This exporter is generated to export the configured transformer initializers")

	// Import TransformerInitializers
	f.ImportAlias("github.com/vulcanize/vulcanizedb/libraries/shared/transformer", "interface")
	for alias, imp := range g.Config.Imports {
		f.ImportAlias(imp, alias)
	}

	// Collect TransformerInitializer names
	importedInitializers := make([]Code, 0, len(g.Config.Imports))
	for _, path := range g.Config.Imports {
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
			"TransformerInitializer").Values(importedInitializers...)))

	// Write code to destination file
	err = f.Save(goFile)
	if err != nil {
		return err
	}

	// Build the .go file into a .so plugin
	return exec.Command("go", "build", "-buildmode=plugin", "-o", soFile, goFile).Run()
}

func GetPaths(config Config) (string, string, error) {
	path, err := homedir.Expand(filepath.Clean(config.FilePath))
	if err != nil {
		return "", "", err
	}
	if strings.Contains(path, "$GOPATH") {
		env := os.Getenv("GOPATH")
		spl := strings.Split(path, "$GOPATH")[1]
		path = filepath.Join(env, spl)
	}

	name := strings.Split(config.FileName, ".")[0]
	goFile := filepath.Join(path, name+".go")
	soFile := filepath.Join(path, name+".so")

	return goFile, soFile, nil
}

func ClearFiles(files ...string) error {
	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			err = os.Remove(file)
			if err != nil {
				return err
			}
		} else if os.IsNotExist(err) {
			// fall through
		} else {
			return err
		}
	}

	return nil
}

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

package manager

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/plugin/helpers"
)

type MigrationManager interface {
	RunMigrations() error
}

type manager struct {
	GenConfig config.Plugin
	DBConfig  config.Database
	tmpMigDir string
}

func NewMigrationManager(gc config.Plugin, dbc config.Database) *manager {
	return &manager{
		GenConfig: gc,
		DBConfig:  dbc,
	}
}

func (m *manager) RunMigrations() error {
	// Get paths to db migrations
	paths, err := m.GenConfig.GetMigrationsPaths()
	if err != nil {
		return err
	}
	if len(paths) < 1 {
		return nil
	}

	// Init directory for temporary copies
	err = m.setupMigrationEnv()
	if err != nil {
		return err
	}
	defer m.cleanUp()

	// Create temporary copies of migrations to the temporary migrationDir
	// These tmps are identical except they have had `1` added in front of their unix_timestamps
	// As such, they will be ran on top of all core migrations (at least, for the next ~317 years)
	// But will still be ran in the same order relative to one another
	// TODO: Less hacky way of handing migrations
	err = m.createMigrationCopies(paths)
	if err != nil {
		return err
	}

	// Run the copied migrations
	pgStr := fmt.Sprintf("postgres://%s:%d/%s?sslmode=disable", m.DBConfig.Hostname, m.DBConfig.Port, m.DBConfig.Name)
	err = exec.Command("migrate", "-path", m.tmpMigDir, "-database", pgStr, "up").Run()
	if err != nil {
		return errors.New(fmt.Sprintf("db migrations for plugin transformers failed: %s", err.Error()))
	}

	return nil
}

func (m *manager) setupMigrationEnv() error {
	// Initialize temp directory for transformer migrations
	var err error
	m.tmpMigDir, err = helpers.CleanPath("$GOPATH/src/github.com/vulcanize/vulcanizedb/db/plugin_migrations")
	if err != nil {
		return err
	}
	err = os.RemoveAll(m.tmpMigDir)
	if err != nil {
		return errors.New(fmt.Sprintf("unable to remove file found at %s where tmp directory needs to be written", m.tmpMigDir))
	}
	err = os.Mkdir(m.tmpMigDir, os.FileMode(0777))
	if err != nil {
		return errors.New(fmt.Sprintf("unable to create temporary migration directory %s", m.tmpMigDir))
	}

	return nil
}

func (m *manager) createMigrationCopies(paths []string) error {
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
			dst := filepath.Join(m.tmpMigDir, "1"+file.Name())
			err = helpers.CopyFile(src, dst)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *manager) cleanUp() error {
	return os.RemoveAll(m.tmpMigDir)
}

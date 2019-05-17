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

package manager

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/lib/pq"
	"github.com/pressly/goose"

	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/plugin/helpers"
)

// Interface for managing the db migrations for plugin transformers
type MigrationManager interface {
	RunMigrations() error
}

type manager struct {
	GenConfig config.Plugin
	DBConfig  config.Database
	tmpMigDir string
	db        *sql.DB
}

// Manager requires both filled in generator and database configs
func NewMigrationManager(gc config.Plugin, dbc config.Database) *manager {
	return &manager{
		GenConfig: gc,
		DBConfig:  dbc,
	}
}

func (m *manager) setDB() error {
	var pgStr string
	if len(m.DBConfig.User) > 0 && len(m.DBConfig.Password) > 0 {
		pgStr = fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
			m.DBConfig.User, m.DBConfig.Password, m.DBConfig.Hostname, m.DBConfig.Port, m.DBConfig.Name)
	} else {
		pgStr = fmt.Sprintf("postgres://%s:%d/%s?sslmode=disable", m.DBConfig.Hostname, m.DBConfig.Port, m.DBConfig.Name)
	}
	dbConnector, err := pq.NewConnector(pgStr)
	if err != nil {
		return errors.New(fmt.Sprintf("can't connect to db: %s", err.Error()))
	}
	m.db = sql.OpenDB(dbConnector)
	return nil
}

func (m *manager) RunMigrations() error {
	// Get paths to db migrations from the plugin config
	paths, err := m.GenConfig.GetMigrationsPaths()
	if err != nil {
		return err
	}
	if len(paths) < 1 {
		return nil
	}
	// Init directory for temporary copies of migrations
	err = m.setupMigrationEnv()
	if err != nil {
		return err
	}
	defer m.cleanUp()
	// Creates copies of migrations for all the plugin's transformers in a tmp dir
	err = m.createMigrationCopies(paths)
	if err != nil {
		return err
	}

	return nil
}

// Setup a temporary directory to hold transformer db migrations
func (m *manager) setupMigrationEnv() error {
	var err error
	m.tmpMigDir, err = helpers.CleanPath(filepath.Join("$GOPATH/src", m.GenConfig.Home, ".plugin_migrations"))
	if err != nil {
		return err
	}
	removeErr := os.RemoveAll(m.tmpMigDir)
	if removeErr != nil {
		removeErrString := "unable to remove file found at %s where tmp directory needs to be written: %s"
		return errors.New(fmt.Sprintf(removeErrString, m.tmpMigDir, removeErr.Error()))
	}
	mkdirErr := os.Mkdir(m.tmpMigDir, os.FileMode(os.ModePerm))
	if mkdirErr != nil {
		mkdirErrString := "unable to create temporary migration directory %s: %s"
		return errors.New(fmt.Sprintf(mkdirErrString, m.tmpMigDir, mkdirErr.Error()))
	}

	return nil
}

// Create copies of db migrations from vendored libs
func (m *manager) createMigrationCopies(paths []string) error {
	// Iterate through migration paths to find migration directory
	for _, path := range paths {
		dir, err := ioutil.ReadDir(path)
		if err != nil {
			return err
		}
		// For each file in the directory check if it is a migration
		for _, file := range dir {
			if file.IsDir() || filepath.Ext(file.Name()) != ".sql" {
				continue
			}
			src := filepath.Join(path, file.Name())
			dst := filepath.Join(m.tmpMigDir, file.Name())
			//  and if it is make a copy of it to our tmp migration directory
			err = helpers.CopyFile(src, dst)
			if err != nil {
				return err
			}
		}
		err = m.fixAndRun(path)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *manager) fixAndRun(path string) error {
	// Setup DB if not set
	if m.db == nil {
		setErr := m.setDB()
		if setErr != nil {
			return errors.New(fmt.Sprintf("could not open db: %s", setErr.Error()))
		}
	}
	// Fix the migrations
	fixErr := goose.Fix(m.tmpMigDir)
	if fixErr != nil {
		return errors.New(fmt.Sprintf("version fixing for plugin migrations at %s failed: %s", path, fixErr.Error()))
	}
	// Run the copied migrations with goose
	upErr := goose.Up(m.db, m.tmpMigDir)
	if upErr != nil {
		return errors.New(fmt.Sprintf("db migrations for plugin transformers at %s failed: %s", path, upErr.Error()))
	}
	return nil
}

func (m *manager) cleanUp() error {
	return os.RemoveAll(m.tmpMigDir)
}

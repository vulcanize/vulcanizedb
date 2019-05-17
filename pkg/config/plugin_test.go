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

package config_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/config"
)

var allDifferentPathsConfig = config.Plugin{
	Transformers: map[string]config.Transformer{
		"transformer1": {
			Path:           "test/init/path",
			Type:           config.EthEvent,
			MigrationPath:  "test/migration/path1",
			MigrationRank:  0,
			RepositoryPath: "test/repo/path",
		},
		"transformer2": {
			Path:           "test/init/path",
			Type:           config.EthEvent,
			MigrationPath:  "test/migration/path2",
			MigrationRank:  2,
			RepositoryPath: "test/repo/path",
		},
		"transformer3": {
			Path:           "test/init/path2",
			Type:           config.EthEvent,
			MigrationPath:  "test/migration/path3",
			MigrationRank:  1,
			RepositoryPath: "test/repo/path",
		},
	},
}

var overlappingPathsConfig = config.Plugin{
	Transformers: map[string]config.Transformer{
		"transformer1": {
			Path:           "test/init/path",
			Type:           config.EthEvent,
			MigrationPath:  "test/migration/path1",
			MigrationRank:  0,
			RepositoryPath: "test/repo/path",
		},
		"transformer2": {
			Path:           "test/init/path",
			Type:           config.EthEvent,
			MigrationPath:  "test/migration/path1",
			MigrationRank:  0,
			RepositoryPath: "test/repo/path",
		},
		"transformer3": {
			Path:           "test/init/path2",
			Type:           config.EthEvent,
			MigrationPath:  "test/migration/path3",
			MigrationRank:  1,
			RepositoryPath: "test/repo/path",
		},
	},
}

var conflictErrorConfig = config.Plugin{
	Transformers: map[string]config.Transformer{
		"transformer1": {
			Path:           "test/init/path",
			Type:           config.EthEvent,
			MigrationPath:  "test/migration/path1",
			MigrationRank:  0,
			RepositoryPath: "test/repo/path",
		},
		"transformer2": {
			Path:           "test/init/path",
			Type:           config.EthEvent,
			MigrationPath:  "test/migration/path2",
			MigrationRank:  0,
			RepositoryPath: "test/repo/path",
		},
		"transformer3": {
			Path:           "test/init/path2",
			Type:           config.EthEvent,
			MigrationPath:  "test/migration/path3",
			MigrationRank:  1,
			RepositoryPath: "test/repo/path",
		},
	},
}

var gapErrorConfig = config.Plugin{
	Transformers: map[string]config.Transformer{
		"transformer1": {
			Path:           "test/init/path",
			Type:           config.EthEvent,
			MigrationPath:  "test/migration/path1",
			MigrationRank:  0,
			RepositoryPath: "test/repo/path",
		},
		"transformer2": {
			Path:           "test/init/path",
			Type:           config.EthEvent,
			MigrationPath:  "test/migration/path2",
			MigrationRank:  3,
			RepositoryPath: "test/repo/path",
		},
		"transformer3": {
			Path:           "test/init/path2",
			Type:           config.EthEvent,
			MigrationPath:  "test/migration/path3",
			MigrationRank:  1,
			RepositoryPath: "test/repo/path",
		},
	},
}

var missingRankErrorConfig = config.Plugin{
	Transformers: map[string]config.Transformer{
		"transformer1": {
			Path:           "test/init/path",
			Type:           config.EthEvent,
			MigrationPath:  "test/migration/path1",
			MigrationRank:  0,
			RepositoryPath: "test/repo/path",
		},
		"transformer2": {
			Path:           "test/init/path",
			Type:           config.EthEvent,
			MigrationPath:  "test/migration/path2",
			RepositoryPath: "test/repo/path",
		},
		"transformer3": {
			Path:           "test/init/path2",
			Type:           config.EthEvent,
			MigrationPath:  "test/migration/path3",
			MigrationRank:  1,
			RepositoryPath: "test/repo/path",
		},
	},
}

var duplicateErrorConfig = config.Plugin{
	Transformers: map[string]config.Transformer{
		"transformer1": {
			Path:           "test/init/path",
			Type:           config.EthEvent,
			MigrationPath:  "test/migration/path1",
			MigrationRank:  0,
			RepositoryPath: "test/repo/path",
		},
		"transformer2": {
			Path:           "test/init/path",
			Type:           config.EthEvent,
			MigrationPath:  "test/migration/path1",
			RepositoryPath: "test/repo/path",
			MigrationRank:  2,
		},
		"transformer3": {
			Path:           "test/init/path2",
			Type:           config.EthEvent,
			MigrationPath:  "test/migration/path3",
			MigrationRank:  1,
			RepositoryPath: "test/repo/path",
		},
	},
}

var _ = Describe("GetMigrationsPaths", func() {
	It("Sorts migration paths by rank", func() {
		plugin := allDifferentPathsConfig
		migrationPaths, err := plugin.GetMigrationsPaths()
		Expect(err).ToNot(HaveOccurred())
		Expect(len(migrationPaths)).To(Equal(3))

		env := os.Getenv("GOPATH")
		path1 := filepath.Join(env, "src/vendor/test/repo/path/test/migration/path1")
		path2 := filepath.Join(env, "src/vendor/test/repo/path/test/migration/path3")
		path3 := filepath.Join(env, "src/vendor/test/repo/path/test/migration/path2")
		expectedMigrationPaths := []string{path1, path2, path3}
		Expect(migrationPaths).To(Equal(expectedMigrationPaths))
	})

	It("Expects identical migration paths to have the same rank", func() {
		plugin := overlappingPathsConfig
		migrationPaths, err := plugin.GetMigrationsPaths()
		Expect(err).ToNot(HaveOccurred())
		Expect(len(migrationPaths)).To(Equal(2))

		env := os.Getenv("GOPATH")
		path1 := filepath.Join(env, "src/vendor/test/repo/path/test/migration/path1")
		path2 := filepath.Join(env, "src/vendor/test/repo/path/test/migration/path3")
		expectedMigrationPaths := []string{path1, path2}
		Expect(migrationPaths).To(Equal(expectedMigrationPaths))
	})

	It("Fails if two different migration paths have the same rank", func() {
		plugin := conflictErrorConfig
		migrationPaths, err := plugin.GetMigrationsPaths()
		Expect(err).To(HaveOccurred())
		Expect(len(migrationPaths)).To(Equal(0))
		Expect(err.Error()).To(ContainSubstring("has the same migration rank"))
	})

	It("Fails if there is a gap in the ranks of the migration paths", func() {
		plugin := gapErrorConfig
		migrationPaths, err := plugin.GetMigrationsPaths()
		Expect(err).To(HaveOccurred())
		Expect(len(migrationPaths)).To(Equal(0))
		Expect(err.Error()).To(ContainSubstring("number of distinct ranks does not match number of distinct migration paths"))
	})

	It("Fails if a transformer is missing its rank", func() {
		plugin := missingRankErrorConfig
		migrationPaths, err := plugin.GetMigrationsPaths()
		Expect(err).To(HaveOccurred())
		Expect(len(migrationPaths)).To(Equal(0))
		Expect(err.Error()).To(ContainSubstring("has the same migration rank"))
	})

	It("Fails if the same migration path has more than one rank", func() {
		plugin := duplicateErrorConfig
		migrationPaths, err := plugin.GetMigrationsPaths()
		Expect(err).To(HaveOccurred())
		Expect(len(migrationPaths)).To(Equal(0))
		Expect(err.Error()).To(ContainSubstring("duplicate paths with different ranks present"))
	})
})

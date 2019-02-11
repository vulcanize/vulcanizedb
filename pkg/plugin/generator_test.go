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

package plugin_test

import (
	"plugin"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"

	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/libraries/shared/watcher"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fs"
	p2 "github.com/vulcanize/vulcanizedb/pkg/plugin"
	"github.com/vulcanize/vulcanizedb/pkg/plugin/helpers"
	"github.com/vulcanize/vulcanizedb/pkg/plugin/test_helpers"
)

var genConfig = config.Plugin{
	Initializers: map[string]string{
		"bite": "github.com/vulcanize/mcd_transformers/transformers/bite/initializer",
		"deal": "github.com/vulcanize/mcd_transformers/transformers/deal/initializer",
	},
	Types: map[string]config.PluginType{
		"bite": config.EthEvent,
		"deal": config.EthEvent,
	},
	Dependencies: map[string]string{
		"mcd_transformers": "github.com/vulcanize/mcd_transformers",
	},
	Migrations: map[string]string{"mcd_transformers": "db/migrations"},
	FileName:   "testEventTransformerSet",
	FilePath:   "$GOPATH/src/github.com/vulcanize/vulcanizedb/pkg/plugin/test_helpers/test",
	Save:       false,
}

var genStorageConfig = config.Plugin{
	Initializers: map[string]string{
		"pit": "github.com/vulcanize/mcd_transformers/transformers/storage_diffs/maker/pit/initializer",
		"vat": "github.com/vulcanize/mcd_transformers/transformers/storage_diffs/maker/vat/initializer",
	},
	Types: map[string]config.PluginType{
		"pit": config.EthStorage,
		"vat": config.EthStorage,
	},
	Dependencies: map[string]string{
		"mcd_transformers": "github.com/vulcanize/mcd_transformers",
	},
	Migrations: map[string]string{"mcd_transformers": "db/migrations"},
	FileName:   "testStorageTransformerSet",
	FilePath:   "$GOPATH/src/github.com/vulcanize/vulcanizedb/pkg/plugin/test_helpers/test",
	Save:       false,
}

var combinedConfig = config.Plugin{
	Initializers: map[string]string{
		"bite": "github.com/vulcanize/mcd_transformers/transformers/bite/initializer",
		"deal": "github.com/vulcanize/mcd_transformers/transformers/deal/initializer",
		"pit":  "github.com/vulcanize/mcd_transformers/transformers/storage_diffs/maker/pit/initializer",
		"vat":  "github.com/vulcanize/mcd_transformers/transformers/storage_diffs/maker/vat/initializer",
	},
	Types: map[string]config.PluginType{
		"bite": config.EthEvent,
		"deal": config.EthEvent,
		"pit":  config.EthStorage,
		"vat":  config.EthStorage,
	},
	Dependencies: map[string]string{
		"mcd_transformers": "github.com/vulcanize/mcd_transformers",
	},
	Migrations: map[string]string{"mcd_transformers": "db/migrations"},
	FileName:   "testComboTransformerSet",
	FilePath:   "$GOPATH/src/github.com/vulcanize/vulcanizedb/pkg/plugin/test_helpers/test",
	Save:       false,
}

var dbConfig = config.Database{
	Hostname: "localhost",
	Port:     5432,
	Name:     "vulcanize_private",
}

type Exporter interface {
	Export() ([]transformer.TransformerInitializer, []transformer.StorageTransformerInitializer)
}

var _ = Describe("Generator test", func() {
	var g p2.Generator
	var goPath, soPath string
	var err error
	var bc core.BlockChain
	var db *postgres.DB
	var hr repositories.HeaderRepository
	var headerID int64
	viper.SetConfigName("compose")
	viper.AddConfigPath("$GOPATH/src/github.com/vulcanize/vulcanizedb/environments/")

	Describe("Event Transformers only", func() {
		BeforeEach(func() {
			goPath, soPath, err = genConfig.GetPluginPaths()
			Expect(err).ToNot(HaveOccurred())
			g, err = p2.NewGenerator(genConfig, dbConfig)
			Expect(err).ToNot(HaveOccurred())
			err = g.GenerateExporterPlugin()
			Expect(err).ToNot(HaveOccurred())
		})
		AfterEach(func() {
			err := helpers.ClearFiles(goPath, soPath)
			Expect(err).ToNot(HaveOccurred())
		})

		Describe("GenerateTransformerPlugin", func() {
			It("It bundles the specified  TransformerInitializers into a Exporter object and creates .so", func() {
				plug, err := plugin.Open(soPath)
				Expect(err).ToNot(HaveOccurred())
				symExporter, err := plug.Lookup("Exporter")
				Expect(err).ToNot(HaveOccurred())
				exporter, ok := symExporter.(Exporter)
				Expect(ok).To(Equal(true))
				initializers, store := exporter.Export()
				Expect(len(initializers)).To(Equal(2))
				Expect(len(store)).To(Equal(0))
			})

			It("Loads our generated Exporter and uses it to import an arbitrary set of TransformerInitializers that we can execute over", func() {
				db, bc = test_helpers.SetupDBandBC()
				defer test_helpers.TearDown(db)

				hr = repositories.NewHeaderRepository(db)
				header1, err := bc.GetHeaderByNumber(9377319)
				Expect(err).ToNot(HaveOccurred())
				headerID, err = hr.CreateOrUpdateHeader(header1)
				Expect(err).ToNot(HaveOccurred())

				plug, err := plugin.Open(soPath)
				Expect(err).ToNot(HaveOccurred())
				symExporter, err := plug.Lookup("Exporter")
				Expect(err).ToNot(HaveOccurred())
				exporter, ok := symExporter.(Exporter)
				Expect(ok).To(Equal(true))
				initializers, _ := exporter.Export()

				w := watcher.NewEventWatcher(db, bc)
				w.AddTransformers(initializers)
				err = w.Execute(constants.HeaderMissing)
				Expect(err).ToNot(HaveOccurred())

				type model struct {
					Ilk              string
					Urn              string
					Ink              string
					Art              string
					IArt             string
					Tab              string
					NFlip            string
					LogIndex         uint   `db:"log_idx"`
					TransactionIndex uint   `db:"tx_idx"`
					Raw              []byte `db:"raw_log"`
					Id               int64  `db:"id"`
					HeaderId         int64  `db:"header_id"`
				}

				returned := model{}

				err = db.Get(&returned, `SELECT * FROM maker.bite WHERE header_id = $1`, headerID)
				Expect(err).ToNot(HaveOccurred())
				Expect(returned.Ilk).To(Equal("4554480000000000000000000000000000000000000000000000000000000000"))
				Expect(returned.Urn).To(Equal("0000000000000000000000000000d8b4147eda80fec7122ae16da2479cbd7ffb"))
				Expect(returned.Ink).To(Equal("80000000000000000000"))
				Expect(returned.Art).To(Equal("11000000000000000000000"))
				Expect(returned.IArt).To(Equal("12496609999999999999992"))
				Expect(returned.Tab).To(Equal("11000000000000000000000"))
				Expect(returned.NFlip).To(Equal("7"))
				Expect(returned.TransactionIndex).To(Equal(uint(1)))
				Expect(returned.LogIndex).To(Equal(uint(4)))
			})
		})
	})

	Describe("Storage Transformers only", func() {
		BeforeEach(func() {
			goPath, soPath, err = genStorageConfig.GetPluginPaths()
			Expect(err).ToNot(HaveOccurred())
			g, err = p2.NewGenerator(genStorageConfig, dbConfig)
			Expect(err).ToNot(HaveOccurred())
			err = g.GenerateExporterPlugin()
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			err := helpers.ClearFiles(goPath, soPath)
			Expect(err).ToNot(HaveOccurred())
		})
		Describe("GenerateTransformerPlugin", func() {
			It("It bundles the specified  StorageTransformerInitializers into a Exporter object and creates .so", func() {
				plug, err := plugin.Open(soPath)
				Expect(err).ToNot(HaveOccurred())
				symExporter, err := plug.Lookup("Exporter")
				Expect(err).ToNot(HaveOccurred())
				exporter, ok := symExporter.(Exporter)
				Expect(ok).To(Equal(true))
				event, initializers := exporter.Export()
				Expect(len(initializers)).To(Equal(2))
				Expect(len(event)).To(Equal(0))
			})

			It("Loads our generated Exporter and uses it to import an arbitrary set of StorageTransformerInitializers that we can execute over", func() {
				db, _ = test_helpers.SetupDBandBC()
				defer test_helpers.TearDown(db)

				plug, err := plugin.Open(soPath)
				Expect(err).ToNot(HaveOccurred())
				symExporter, err := plug.Lookup("Exporter")
				Expect(err).ToNot(HaveOccurred())
				exporter, ok := symExporter.(Exporter)
				Expect(ok).To(Equal(true))
				_, initializers := exporter.Export()

				tailer := fs.FileTailer{Path: viper.GetString("filesystem.storageDiffsPath")}
				w := watcher.NewStorageWatcher(tailer, db)
				w.AddTransformers(initializers)
				// This blocks right now, need to make test file to read from
				//err = w.Execute()
				//Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("Event and Storage Transformers in same instance", func() {
		BeforeEach(func() {
			goPath, soPath, err = combinedConfig.GetPluginPaths()
			Expect(err).ToNot(HaveOccurred())
			g, err = p2.NewGenerator(combinedConfig, dbConfig)
			Expect(err).ToNot(HaveOccurred())
			err = g.GenerateExporterPlugin()
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			err := helpers.ClearFiles(goPath, soPath)
			Expect(err).ToNot(HaveOccurred())
		})
		Describe("GenerateTransformerPlugin", func() {
			It("It bundles the specified  TransformerInitializers and StorageTransformerInitializers into a Exporter object and creates .so", func() {
				plug, err := plugin.Open(soPath)
				Expect(err).ToNot(HaveOccurred())
				symExporter, err := plug.Lookup("Exporter")
				Expect(err).ToNot(HaveOccurred())
				exporter, ok := symExporter.(Exporter)
				Expect(ok).To(Equal(true))
				eventInitializers, storageInitializers := exporter.Export()
				Expect(len(eventInitializers)).To(Equal(2))
				Expect(len(storageInitializers)).To(Equal(2))
			})

			It("Loads our generated Exporter and uses it to import an arbitrary set of TransformerInitializers and StorageTransformerInitializers that we can execute over", func() {
				db, bc = test_helpers.SetupDBandBC()
				defer test_helpers.TearDown(db)

				hr = repositories.NewHeaderRepository(db)
				header1, err := bc.GetHeaderByNumber(9377319)
				Expect(err).ToNot(HaveOccurred())
				headerID, err = hr.CreateOrUpdateHeader(header1)
				Expect(err).ToNot(HaveOccurred())

				plug, err := plugin.Open(soPath)
				Expect(err).ToNot(HaveOccurred())
				symExporter, err := plug.Lookup("Exporter")
				Expect(err).ToNot(HaveOccurred())
				exporter, ok := symExporter.(Exporter)
				Expect(ok).To(Equal(true))
				eventInitializers, storageInitializers := exporter.Export()

				ew := watcher.NewEventWatcher(db, bc)
				ew.AddTransformers(eventInitializers)
				err = ew.Execute(constants.HeaderMissing)
				Expect(err).ToNot(HaveOccurred())

				type model struct {
					Ilk              string
					Urn              string
					Ink              string
					Art              string
					IArt             string
					Tab              string
					NFlip            string
					LogIndex         uint   `db:"log_idx"`
					TransactionIndex uint   `db:"tx_idx"`
					Raw              []byte `db:"raw_log"`
					Id               int64  `db:"id"`
					HeaderId         int64  `db:"header_id"`
				}

				returned := model{}

				err = db.Get(&returned, `SELECT * FROM maker.bite WHERE header_id = $1`, headerID)
				Expect(err).ToNot(HaveOccurred())
				Expect(returned.Ilk).To(Equal("4554480000000000000000000000000000000000000000000000000000000000"))
				Expect(returned.Urn).To(Equal("0000000000000000000000000000d8b4147eda80fec7122ae16da2479cbd7ffb"))
				Expect(returned.Ink).To(Equal("80000000000000000000"))
				Expect(returned.Art).To(Equal("11000000000000000000000"))
				Expect(returned.IArt).To(Equal("12496609999999999999992"))
				Expect(returned.Tab).To(Equal("11000000000000000000000"))
				Expect(returned.NFlip).To(Equal("7"))
				Expect(returned.TransactionIndex).To(Equal(uint(1)))
				Expect(returned.LogIndex).To(Equal(uint(4)))

				tailer := fs.FileTailer{Path: viper.GetString("filesystem.storageDiffsPath")}
				sw := watcher.NewStorageWatcher(tailer, db)
				sw.AddTransformers(storageInitializers)
				// This blocks right now, need to make test file to read from
				//err = w.Execute()
				//Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})

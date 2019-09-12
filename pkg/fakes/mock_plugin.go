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

package fakes

import (
	"github.com/vulcanize/vulcanizedb/libraries/shared/mocks"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"plugin"
)

type MockPlugin struct {
	FakeEventInitializers    []transformer.EventTransformerInitializer
	FakeStorageInitializers  []transformer.StorageTransformerInitializer
	FakeContractInitializers []transformer.ContractTransformerInitializer
}

type exporter string

var Exporter exporter

func (e exporter) Export() ([]transformer.EventTransformerInitializer, []transformer.StorageTransformerInitializer, []transformer.ContractTransformerInitializer) {
	fakeEventTransformer := &mocks.MockTransformer{}
	fakeStorageTransformer := &mocks.MockStorageTransformer{}
	eventTransformerInitializers := []transformer.EventTransformerInitializer{fakeEventTransformer.FakeTransformerInitializer}
	storageTransformerInitializers := []transformer.StorageTransformerInitializer{fakeStorageTransformer.FakeTransformerInitializer}
	return eventTransformerInitializers, storageTransformerInitializers, []transformer.ContractTransformerInitializer{}
}

func NewMockPlugin() MockPlugin {
	e, s, c := Exporter.Export()
	return MockPlugin{
		FakeEventInitializers:    e,
		FakeStorageInitializers:  s,
		FakeContractInitializers: c,
	}
}

func (p MockPlugin) Lookup(symName string) (plugin.Symbol, error) {
	return Exporter, nil
}

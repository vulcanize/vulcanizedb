// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package event_triggered

import (
	"log"

	"fmt"
	"github.com/vulcanize/vulcanizedb/examples/constants"
	"github.com/vulcanize/vulcanizedb/examples/generic"
	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
)

type DaiTransformer struct {
	Converter              ERC20ConverterInterface
	WatchedEventRepository datastore.WatchedEventRepository
	FilterRepository       datastore.FilterRepository
	Repository             Datastore
}

func NewTransformer(db *postgres.DB, config generic.ContractConfig) shared.Transformer {
	var transformer shared.Transformer
	cnvtr := NewERC20Converter(config)
	wer := repositories.WatchedEventRepository{DB: db}
	fr := repositories.FilterRepository{DB: db}
	lkr := Repository{DB: db}
	transformer = &DaiTransformer{
		Converter:              cnvtr,
		WatchedEventRepository: wer,
		FilterRepository:       fr,
		Repository:             lkr,
	}
	for _, filter := range constants.DaiFilters {
		fr.CreateFilter(filter)
	}
	return transformer
}

func (tr DaiTransformer) Execute() error {
	for _, filter := range constants.DaiFilters {
		watchedEvents, err := tr.WatchedEventRepository.GetWatchedEvents(filter.Name)
		if err != nil {
			log.Println(fmt.Sprintf("Error fetching events for %s:", filter.Name), err)
			return err
		}
		for _, we := range watchedEvents {
			if filter.Name == "Transfer" {
				entity, err := tr.Converter.ToTransferEntity(*we)
				model := tr.Converter.ToTransferModel(*entity)
				if err != nil {
					log.Printf("Error persisting data for Dai Transfers (watchedEvent.LogID %d):\n %s", we.LogID, err)
				}
				tr.Repository.CreateTransfer(model, we.LogID)
			}
			if filter.Name == "Approval" {
				entity, err := tr.Converter.ToApprovalEntity(*we)
				model := tr.Converter.ToApprovalModel(*entity)
				if err != nil {
					log.Printf("Error persisting data for Dai Approvals (watchedEvent.LogID %d):\n %s", we.LogID, err)
				}
				tr.Repository.CreateApproval(model, we.LogID)
			}
		}
	}
	return nil
}

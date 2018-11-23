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

package dai

import (
	"fmt"
	"log"

	"github.com/vulcanize/vulcanizedb/examples/erc20_watcher/event_triggered"
	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
	"github.com/vulcanize/vulcanizedb/pkg/omni/shared/constants"
)

type ERC20EventTransformer struct {
	Converter              ERC20ConverterInterface
	WatchedEventRepository datastore.WatchedEventRepository
	FilterRepository       datastore.FilterRepository
	Repository             event_triggered.ERC20EventDatastore
	Filters                []filters.LogFilter
}

func NewTransformer(db *postgres.DB, blockchain core.BlockChain, con shared.ContractConfig) (shared.Transformer, error) {
	var transformer shared.Transformer

	cnvtr, err := NewERC20Converter(con)
	if err != nil {
		return nil, err
	}

	wer := repositories.WatchedEventRepository{DB: db}
	fr := repositories.FilterRepository{DB: db}
	lkr := event_triggered.ERC20EventRepository{DB: db}

	transformer = ERC20EventTransformer{
		Converter:              cnvtr,
		WatchedEventRepository: wer,
		FilterRepository:       fr,
		Repository:             lkr,
		Filters:                con.Filters,
	}

	for _, filter := range con.Filters {
		fr.CreateFilter(filter)
	}

	return transformer, nil
}

func (tr ERC20EventTransformer) Execute() error {
	for _, filter := range tr.Filters {
		watchedEvents, err := tr.WatchedEventRepository.GetWatchedEvents(filter.Name)
		if err != nil {
			log.Println(fmt.Sprintf("Error fetching events for %s:", filter.Name), err)
			return err
		}
		for _, we := range watchedEvents {
			if filter.Name == constants.TransferEvent.String() {
				entity, err := tr.Converter.ToTransferEntity(*we)
				model := tr.Converter.ToTransferModel(entity)
				if err != nil {
					log.Printf("Error persisting data for Dai Transfers (watchedEvent.LogID %d):\n %s", we.LogID, err)
				}
				tr.Repository.CreateTransfer(model, we.LogID)
			}
			if filter.Name == constants.ApprovalEvent.String() {
				entity, err := tr.Converter.ToApprovalEntity(*we)
				model := tr.Converter.ToApprovalModel(entity)
				if err != nil {
					log.Printf("Error persisting data for Dai Approvals (watchedEvent.LogID %d):\n %s", we.LogID, err)
				}
				tr.Repository.CreateApproval(model, we.LogID)
			}
		}
	}

	return nil
}

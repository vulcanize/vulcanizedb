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

import "github.com/vulcanize/vulcanizedb/pkg/core"

type MockHeaderSyncHeaderRepository struct{}

func (*MockHeaderSyncHeaderRepository) AddCheckColumn(id string) error {
	return nil
}

func (*MockHeaderSyncHeaderRepository) AddCheckColumns(ids []string) error {
	panic("implement me")
}

func (*MockHeaderSyncHeaderRepository) MarkHeaderChecked(headerID int64, eventID string) error {
	panic("implement me")
}

func (*MockHeaderSyncHeaderRepository) MarkHeaderCheckedForAll(headerID int64, ids []string) error {
	panic("implement me")
}

func (*MockHeaderSyncHeaderRepository) MarkHeadersCheckedForAll(headers []core.Header, ids []string) error {
	panic("implement me")
}

func (*MockHeaderSyncHeaderRepository) MissingHeaders(startingBlockNumber int64, endingBlockNumber int64, eventID string) ([]core.Header, error) {
	panic("implement me")
}

func (*MockHeaderSyncHeaderRepository) MissingMethodsCheckedEventsIntersection(startingBlockNumber, endingBlockNumber int64, methodIds, eventIds []string) ([]core.Header, error) {
	panic("implement me")
}

func (*MockHeaderSyncHeaderRepository) MissingHeadersForAll(startingBlockNumber, endingBlockNumber int64, ids []string) ([]core.Header, error) {
	panic("implement me")
}

func (*MockHeaderSyncHeaderRepository) CheckCache(key string) (interface{}, bool) {
	panic("implement me")
}

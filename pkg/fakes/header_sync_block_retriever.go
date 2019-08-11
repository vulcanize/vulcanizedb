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

type MockHeaderSyncBlockRetriever struct {
	FirstBlock    int64
	FirstBlockErr error
}

func (retriever *MockHeaderSyncBlockRetriever) RetrieveFirstBlock() (int64, error) {
	return retriever.FirstBlock, retriever.FirstBlockErr
}

func (retriever *MockHeaderSyncBlockRetriever) RetrieveMostRecentBlock() (int64, error) {
	return 0, nil
}

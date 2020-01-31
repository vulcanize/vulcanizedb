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

package shared

// These types serve as very loose wrappers around a generic underlying interface{}
type RawChainData interface{}

// The concrete type underneath StreamedIPLDs can be a pointer only if the Value() method returns a copy of the values
// stored at that memory location and not a copy of the pointer itself.
// We want to avoid sending a pointer to publishAndIndex and screenAndServe channels; sharing memory across these processes
type StreamedIPLDs interface {
	Value() StreamedIPLDs
}

type CIDsForIndexing interface{}

type CIDsForFetching interface{}

type FetchedIPLDs interface{}

// The concrete type underneath StreamedIPLDs can be a pointer only if the Value() method returns a copy of the values
// stored at that memory location and not a copy of the pointer itself.
// We want to avoid sending a pointer to subscription channels; sharing memory across all subscriptions
type ServerResponse interface {
	Value() ServerResponse
}

type Gap struct {
	Start uint64
	Stop  uint64
}

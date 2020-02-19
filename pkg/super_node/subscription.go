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

package super_node

import (
	"errors"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/super_node/shared"
)

type Flag int32

const (
	EmptyFlag Flag = iota
	BackFillCompleteFlag
)

// Subscription holds the information for an individual client subscription to the super node
type Subscription struct {
	ID          rpc.ID
	PayloadChan chan<- SubscriptionPayload
	QuitChan    chan<- bool
}

// SubscriptionPayload is the struct for a super node stream payload
// It carries data of a type specific to the chain being supported/queried and an error message
type SubscriptionPayload struct {
	Data shared.ServerResponse `json:"data"` // e.g. for Ethereum eth.StreamPayload
	Err  string                `json:"err"`  // field for error
	Flag Flag                  `json:"flag"` // field for message
}

func (sp SubscriptionPayload) Error() error {
	if sp.Err == "" {
		return nil
	}
	return errors.New(sp.Err)
}

func (sp SubscriptionPayload) BackFillComplete() bool {
	if sp.Flag == BackFillCompleteFlag {
		return true
	}
	return false
}

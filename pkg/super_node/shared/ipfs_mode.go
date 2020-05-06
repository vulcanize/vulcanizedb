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

import (
	"errors"
	"strings"
)

// IPFSMode enum for specifying how we want to interface and publish objects to IPFS
type IPFSMode int

const (
	Unknown IPFSMode = iota
	LocalInterface
	RemoteClient
	DirectPostgres
)

func (c IPFSMode) String() string {
	switch c {
	case LocalInterface:
		return "Local"
	case RemoteClient:
		return "Remote"
	case DirectPostgres:
		return "Postgres"
	default:
		return ""
	}
}

func NewIPFSMode(name string) (IPFSMode, error) {
	switch strings.ToLower(name) {
	case "local", "interface", "minimal":
		return LocalInterface, nil
	case "remote", "client":
		return RemoteClient, errors.New("remote IPFS client mode is not currently supported")
	case "postgres", "direct":
		return DirectPostgres, nil
	default:
		return Unknown, errors.New("unrecognized name for ipfs mode")
	}
}

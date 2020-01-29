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
	"github.com/hpcloud/tail"
	"gopkg.in/tomb.v1"
)

type MockTailer struct {
	Lines   chan *tail.Line
	TailErr error
}

func NewMockTailer() *MockTailer {
	return &MockTailer{
		Lines: make(chan *tail.Line, 1),
	}
}

func (mock *MockTailer) Tail() (*tail.Tail, error) {
	fakeTail := &tail.Tail{
		Filename: "",
		Lines:    mock.Lines,
		Config:   tail.Config{},
		Tomb:     tomb.Tomb{},
	}
	return fakeTail, mock.TailErr
}

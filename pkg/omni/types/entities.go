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

package types

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Config struct {
	Network string
	BC      core.BlockChain
	DB      *postgres.DB
}

type Event struct {
	Name      string
	Anonymous bool
	Fields    []*Field
	Logs      map[int64]Log // Map of VulcanizeIdLog to parsed event log
}

type Method struct {
	Name    string
	Const   bool
	Inputs  []*Field
	Outputs []*Field
}

type Field struct {
	abi.Argument
	PgType string
}

type Log struct {
	Values map[string]interface{} // Map of event input names to their values
	Block  int64
	Tx     string
}

func NewEvent(e abi.Event) *Event {
	fields := make([]*Field, len(e.Inputs))
	for i, input := range e.Inputs {
		fields[i] = &Field{}
		fields[i].Name = input.Name
		fields[i].Type = input.Type
		fields[i].Indexed = input.Indexed
	}

	return &Event{
		Name:      e.Name,
		Anonymous: e.Anonymous,
		Fields:    fields,
		Logs:      map[int64]Log{},
	}
}

func NewMethod(m abi.Method) *Method {
	inputs := make([]*Field, len(m.Inputs))
	for i, input := range m.Inputs {
		inputs[i] = &Field{}
		inputs[i].Name = input.Name
		inputs[i].Type = input.Type
		inputs[i].Indexed = input.Indexed
	}

	outputs := make([]*Field, len(m.Outputs))
	for i, output := range m.Outputs {
		outputs[i] = &Field{}
		outputs[i].Name = output.Name
		outputs[i].Type = output.Type
		outputs[i].Indexed = output.Indexed
	}

	return &Method{
		Name:    m.Name,
		Const:   m.Const,
		Inputs:  inputs,
		Outputs: outputs,
	}
}

func (e Event) Sig() string {
	types := make([]string, len(e.Fields))

	for i, input := range e.Fields {
		types[i] = input.Type.String()
	}

	return fmt.Sprintf("%v(%v)", e.Name, strings.Join(types, ","))
}

func (m Method) Sig() string {
	types := make([]string, len(m.Inputs))
	i := 0
	for _, input := range m.Inputs {
		types[i] = input.Type.String()
		i++
	}

	return fmt.Sprintf("%v(%v)", m.Name, strings.Join(types, ","))
}

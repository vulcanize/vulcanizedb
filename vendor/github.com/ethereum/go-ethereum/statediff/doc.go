// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

/*
This work is adapted from work by Charles Crain at https://github.com/jpmorganchase/quorum/blob/9b7fd9af8082795eeeb6863d9746f12b82dd5078/statediff/statediff.go

Package statediff provides an auxiliary service that processes state diff objects from incoming chain events,
relaying the objects to any rpc subscriptions.

Rpc subscriptions to the service can be created using the rpc.Client.Subscribe() method,
with the "statediff" namespace, a statediff.Payload channel, and the name of the statediff api's rpc method- "stream".

e.g.

stateDiffPayloadChan := make(chan statediff.Payload, 20000)
rpcSub, err := Subscribe(context.Background(), "statediff", stateDiffPayloadChan, "stream"})
for {
	select {
	case stateDiffPayload := <- stateDiffPayloadChan:
		processPayload(stateDiffPayload)
	case err := <= rpcSub.Err():
		log.Error(err)
	}
}

*/
package statediff

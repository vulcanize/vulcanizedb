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

import log "github.com/sirupsen/logrus"

func sendNonBlockingErr(sub Subscription, err error) {
	log.Error(err)
	select {
	case sub.PayloadChan <- SubscriptionPayload{nil, err.Error()}:
	default:
		log.Infof("unable to send error to subscription %s", sub.ID)
	}
}

func sendNonBlockingQuit(sub Subscription) {
	select {
	case sub.QuitChan <- true:
		log.Infof("closing subscription %s", sub.ID)
	default:
		log.Infof("unable to close subscription %s; channel has no receiver", sub.ID)
	}
}

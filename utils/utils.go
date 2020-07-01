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

package utils

import (
	"github.com/jmoiron/sqlx"
	"github.com/makerdao/vulcanizedb/pkg/config"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/sirupsen/logrus"
)

func LoadPostgres(database config.Database, node core.Node) postgres.DB {
	db, err := postgres.NewDB(database, node)
	if err != nil {
		logrus.Fatal("Error loading postgres: ", err)
	}
	return *db
}

func RollbackAndLogFailure(tx *sqlx.Tx, txErr error, fieldName string) {
	rollbackErr := tx.Rollback()
	if rollbackErr != nil {
		logrus.WithFields(logrus.Fields{"rollbackErr": rollbackErr, "txErr": txErr}).
			Warnf("failed to rollback transaction after failing to insert %s", fieldName)
	}
}

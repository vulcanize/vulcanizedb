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

package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Env variables
const (
	DATABASE_NAME     = "DATABASE_NAME"
	DATABASE_HOSTNAME = "DATABASE_HOSTNAME"
	DATABASE_PORT     = "DATABASE_PORT"
	DATABASE_USER     = "DATABASE_USER"
	DATABASE_PASSWORD = "DATABASE_PASSWORD"
)

type Database struct {
	Hostname string
	Name     string
	User     string
	Password string
	Port     int
}

func DbConnectionString(dbConfig Database) string {
	if len(dbConfig.User) > 0 && len(dbConfig.Password) > 0 {
		return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
			dbConfig.User, dbConfig.Password, dbConfig.Hostname, dbConfig.Port, dbConfig.Name)
	}
	if len(dbConfig.User) > 0 && len(dbConfig.Password) == 0 {
		return fmt.Sprintf("postgresql://%s@%s:%d/%s?sslmode=disable",
			dbConfig.User, dbConfig.Hostname, dbConfig.Port, dbConfig.Name)
	}
	return fmt.Sprintf("postgresql://%s:%d/%s?sslmode=disable", dbConfig.Hostname, dbConfig.Port, dbConfig.Name)
}

func (d *Database) Init() {
	viper.BindEnv("database.name", DATABASE_NAME)
	viper.BindEnv("database.hostname", DATABASE_HOSTNAME)
	viper.BindEnv("database.port", DATABASE_PORT)
	viper.BindEnv("database.user", DATABASE_USER)
	viper.BindEnv("database.password", DATABASE_PASSWORD)
	d.Name = viper.GetString("database.name")
	d.Hostname = viper.GetString("database.hostname")
	d.Port = viper.GetInt("database.port")
	d.User = viper.GetString("database.user")
	d.Password = viper.GetString("database.password")
}

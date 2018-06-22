package config

import "fmt"

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
	return fmt.Sprintf("postgresql://%s:%d/%s?sslmode=disable", dbConfig.Hostname, dbConfig.Port, dbConfig.Name)
}

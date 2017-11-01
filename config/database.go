package config

import "fmt"

type Database struct {
	Hostname string
	Name     string
	Port     int
}

func DbConnectionString(dbConfig Database) string {
	return fmt.Sprintf("postgresql://%s:%d/%s?sslmode=disable", dbConfig.Hostname, dbConfig.Port, dbConfig.Name)
}

package db

type DatabaseType int

const (
	Level DatabaseType = iota
)

type DatabaseConfig struct {
	Type DatabaseType
	Path string
}

func CreateDatabaseConfig(dbType DatabaseType, path string) DatabaseConfig {
	return DatabaseConfig{
		Type: dbType,
		Path: path,
	}
}

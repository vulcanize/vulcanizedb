package config

func Private() Config {
	return Config{
		Database: Database{
			Name:     "vulcanize_private",
			Hostname: "localhost",
			Port:     5432,
		},
	}
}

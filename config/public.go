package config

func Public() Config {
	return Config{
		Database: Database{
			Name:     "vulcanize_public",
			Hostname: "localhost",
			Port:     5432,
		},
	}
}

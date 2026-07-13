package config

type Config struct {
	Db_url string
	User   string
}

func Read(dbURL string) *Config {
	return &Config{Db_url: dbURL, User: "Fin"}
}

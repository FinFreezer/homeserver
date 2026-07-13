package config

type UserConfig struct {
	Db_url string
	User   string
}

func Read(dbURL, user string) *UserConfig {
	return &UserConfig{Db_url: dbURL, User: user}
}

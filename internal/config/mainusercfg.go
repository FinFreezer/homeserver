package config

type UserConfig struct {
	User       string
	Password   string
	Authorized bool
}

func Read(user, password string) *UserConfig {
	return &UserConfig{User: user, Password: password, Authorized: false}
}

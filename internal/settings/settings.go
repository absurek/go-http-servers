package settings

import "os"

type Settings struct {
	DBUrl     string
	JWTSecret string
	PolkaKey  string
}

func NewSettings() Settings {
	return Settings{
		DBUrl:     os.Getenv("DB_URL"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		PolkaKey:  os.Getenv("POLKA_KEY"),
	}
}

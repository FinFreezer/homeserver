package functionality

import (
	"database/sql"
	"log"
	"os"

	"github.com/finfreezer/homeserver/internal/database"
	"github.com/joho/godotenv"
)

type ApiConfig struct {
	Database   *database.Queries
	Platform   string
	Secret     string
	ApiKey     string
	Authorized bool
}

func OpenDatabase() (*ApiConfig, error) {
	log.Println("Loading server data...")
	godotenv.Load()
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("SECRET")
	dbURL := os.Getenv("DB_URL")
	ApiKey := os.Getenv("API_KEY")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	dbQueries := database.New(db)
	newApiConf := ApiConfig{
		Database:   dbQueries,
		Platform:   platform,
		Secret:     secret,
		ApiKey:     ApiKey,
		Authorized: false,
	}
	return &newApiConf, nil
}

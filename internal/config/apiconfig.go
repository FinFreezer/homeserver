package config

import (
	"database/sql"
	"log"
	"os"
	"sync/atomic"

	"github.com/finfreezer/homeserver/internal/database"
	"github.com/joho/godotenv"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	Database       *database.Queries
	Platform       string
	Secret         string
	ApiKey         string
}

func OpenDatabase(args []string) (*ApiConfig, *UserConfig, error) {
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("SECRET")
	dbURL := os.Getenv("DB_URL")
	ApiKey := os.Getenv("API_KEY")

	if len(args) != 1 {
		log.Println("Loading client data...")
		godotenv.Load()
		newUserConf := Read(dbURL, args[1])

		db, err := sql.Open("postgres", newUserConf.Db_url)
		if err != nil {
			log.Println(err)
			return nil, nil, err
		}
		dbQueries := database.New(db)
		newApiConf := ApiConfig{
			FileserverHits: atomic.Int32{},
			Database:       dbQueries,
			Platform:       platform,
			Secret:         secret,
			ApiKey:         ApiKey,
		}
		return &newApiConf, newUserConf, nil

	} else {
		log.Println("Loading server data...")
		db, err := sql.Open("postgres", dbURL)
		if err != nil {
			log.Println(err)
			return nil, nil, err
		}
		dbQueries := database.New(db)
		newApiConf := ApiConfig{
			FileserverHits: atomic.Int32{},
			Database:       dbQueries,
			Platform:       platform,
			Secret:         secret,
			ApiKey:         ApiKey,
		}
		return &newApiConf, nil, nil
	}
}

package main

import (
	"fmt"
	"log"
	"os"

	"database/sql"

	"github.com/finfreezer/homeserver/internal/auth"
	c "github.com/finfreezer/homeserver/internal/config"
	"github.com/finfreezer/homeserver/internal/database"
	pg "github.com/finfreezer/homeserver/internal/functionality"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	newConf := c.Read(os.Getenv("DB_URL"))
	db, err := sql.Open("postgres", newConf.Db_url)
	args := os.Args
	if len(args) != 3 {
		fmt.Println("Add the admin username and password.")
		os.Exit(1)
	}
	newConf.User = args[1]
	pwHash, err := auth.CreatePasswordHash(args[2])

	if err != nil {
		log.Println(err)
	}
	dbQueries := database.New(db)
	newState := pg.State{Db: dbQueries, Cfg: newConf}
	pg.AddAdmin(&newState, pwHash)
}

package functionality

import (
	"context"
	"log"

	"github.com/finfreezer/homeserver/internal/database"
)

func AddAdmin(db *database.Queries, username, passwordhash string) {
	params := database.CreateUserParams{Name: username, PasswordHash: passwordhash}
	user, err := db.CreateUser(context.Background(), params)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(user.Name)
}

package functionality

import (
	"context"
	"log"

	"github.com/finfreezer/homeserver/internal/auth"
	"github.com/finfreezer/homeserver/internal/database"
)

func AddAdmin(db *database.Queries, username, passwordhash, password string) bool {
	params := database.CreateUserParams{Name: username, PasswordHash: passwordhash}
	log.Println("Looking if user exists.")
	user, err := db.FindUser(context.Background(), params.Name)
	if err == nil {
		match, err := auth.CheckPassword(password, user.PasswordHash)
		if match && err == nil {
			log.Printf("Logged in as admin: %s", user.Name)
			return true
		}
		if !match || err != nil {
			log.Println("Error logging in as user, check password.")
			return false
		}
	}

	log.Println("No user found, creating admin user.")
	user, err = db.CreateUser(context.Background(), params)
	if err != nil {
		log.Fatal(err)
		return false
	}
	log.Printf("Logged in as admin: %s", user.Name)
	return true
}

package functionality

import (
	"context"
	"log"

	c "github.com/finfreezer/homeserver/internal/config"
	"github.com/finfreezer/homeserver/internal/database"
)

type State struct {
	Db  *database.Queries
	Cfg *c.Config
}

func AddAdmin(state *State, passwordhash string) {
	params := database.CreateUserParams{Name: state.Cfg.User, PasswordHash: passwordhash}
	user, err := state.Db.CreateUser(context.Background(), params)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(user.Name)
}

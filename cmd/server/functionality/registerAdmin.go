package functionality

import (
	"context"
	"log"

	c "github.com/finfreezer/homeserver/internal/config"
	"github.com/finfreezer/homeserver/internal/database"
)

type State struct {
	CfgUser *c.UserConfig
	CfgApi  *c.ApiConfig
}

func AddAdmin(state *State, passwordhash string) {
	params := database.CreateUserParams{Name: state.CfgUser.User, PasswordHash: passwordhash}
	user, err := state.CfgApi.Database.CreateUser(context.Background(), params)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(user.Name)
}

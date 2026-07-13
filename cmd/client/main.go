package main

import (
	"log"
	"os"

	"github.com/finfreezer/homeserver/internal/auth"
	c "github.com/finfreezer/homeserver/internal/config"
	pg "github.com/finfreezer/homeserver/internal/functionality"
	_ "github.com/lib/pq"
)

func main() {
	args := os.Args
	newApiConf, newUserConf, err := c.OpenDatabase(args)
	if err != nil {
		log.Fatalf(err.Error())
	}
	newState := pg.State{CfgUser: newUserConf, CfgApi: newApiConf}
	pwHash, err := auth.CreatePasswordHash(args[2])

	if err != nil {
		log.Println(err)
	}
	pg.AddAdmin(&newState, pwHash)
}

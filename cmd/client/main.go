package main

import (
	pg "github.com/finfreezer/homeserver/cmd/client/functionality"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	pg.MainCLI()
}

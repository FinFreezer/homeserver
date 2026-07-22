package main

import (
	pg "github.com/finfreezer/homeserver/cmd/client/functionality"
	//"github.com/pkg/browser"
	"github.com/joho/godotenv"
)

func main() {
	/*const url = "http://golang.org/"
	browser.OpenURL(url)*/
	godotenv.Load()
	pg.MainCLI()
}

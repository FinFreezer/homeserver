package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	pg "github.com/finfreezer/homeserver/cmd/server/functionality"
	"github.com/finfreezer/homeserver/internal/auth"
	_ "github.com/lib/pq"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Please input login credentials.")
		os.Exit(1)
	}
	if args[1] == "delete" {
		newApiConf, err := pg.OpenDatabase()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		newApiConf.Database.DeleteUsers(context.Background())
		fmt.Println("Cleared 'users' table.")
		os.Exit(0)
	}
	if len(args) != 3 {
		fmt.Println("Please input the username and password you wish to use remotely.")
		os.Exit(1)
	}
	newApiConf, err := pg.OpenDatabase()
	if err != nil {
		log.Println(err)
	}
	pwHash, err := auth.CreatePasswordHash(args[2]) //Username + hashed password
	if ok := pg.AddAdmin(newApiConf.Database, args[1], pwHash, args[2]); !ok {
		log.Println("Error logging in to the server.")
		os.Exit(1)
	}
	newMux := http.NewServeMux()
	newMux.HandleFunc("POST /login", newApiConf.Login)
	newMux.HandleFunc("GET /listdir/{path...}", newApiConf.ListContents)
	newMux.HandleFunc("GET /stream/{path...}", newApiConf.StreamVideo)
	newMux.HandleFunc("PUT /cd", newApiConf.MoveRootDirectory)
	//newServer := http.Server{Addr: ":8080", Handler: newMux}
	newServer := http.FileServer(http.Dir("./cmd/server"))
	newMux.Handle("/", newServer)
	//err = newServer.ListenAndServe()
	log.Println("Listening and Serving on port :12000")
	http.ListenAndServe(":12000", newMux)
	if err != nil {
		log.Fatal(err)
	}
}

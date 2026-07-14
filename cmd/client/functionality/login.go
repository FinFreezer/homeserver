package functionality

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	c "github.com/finfreezer/homeserver/internal/config"
	"github.com/joho/godotenv"
)

type loginParameters struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type responseParams struct {
	Message string `json:"reply"`
	Error   string `json:"error"`
}

func login(cfg *c.UserConfig) bool {
	godotenv.Load()
	params := loginParameters{cfg.User, cfg.Password}
	rqst, err := json.Marshal(params)
	if err != nil {
		log.Println(err)
		return false
	}
	fullUrl := os.Getenv("DST_SERVER") + "/login"
	log.Printf("Making a POST request to %s...\n", fullUrl)
	if err != nil {
		log.Println(err)
		return false
	}
	resp, err := http.Post(fullUrl, "application/json", bytes.NewReader(rqst))
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	reply := responseParams{}
	err = decoder.Decode(&reply)
	log.Printf("%+v", reply)
	return true
}

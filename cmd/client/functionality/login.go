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
	log.Printf("Making a POST request to %s\n", fullUrl)
	resp, err := http.Post(fullUrl, "application/json", bytes.NewReader(rqst))
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	reply := responseParams{}
	if err = decoder.Decode(&reply); err != nil {
		log.Println(err)
		return false
	}
	if 200 <= resp.StatusCode && resp.StatusCode < 300 {
		log.Printf("%s", reply.Message)
		cfg.Authorized = true
		return true
	} else {
		log.Printf("%s", reply.Error)
		cfg.Authorized = false
		return false
	}
}

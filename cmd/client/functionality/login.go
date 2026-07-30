package functionality

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	c "github.com/finfreezer/homeserver/internal/config"
	"github.com/joho/godotenv"
)

type loginParameters struct {
	Name      string `json:"name"`
	Password  string `json:"password"`
	WithToken bool   `json:"withToken"`
	Token     string `json:"token,omitempty"`
}

type responseParams struct {
	Message string `json:"reply"`
	Token   string `json:"token,omitempty"`
	Error   string `json:"error,omitempty"`
}

// Todo: refactor to using client with newRequest and .Do()
func login(cfg *c.UserConfig) bool {
	params := loginParameters{Name: cfg.User, WithToken: true, Token: cfg.Token}
	if len(cfg.Args) > 1 {
		cfg.User = cfg.Args[0]
		cfg.Password = cfg.Args[1]
		params = loginParameters{Name: cfg.User, WithToken: false, Password: cfg.Password, Token: ""}
	}
	rqst, err := json.Marshal(params)
	if err != nil {
		log.Println(err)
		return false
	}
	fullUrl := os.Getenv("DST_SERVER") + os.Getenv("DFLT_PORT") + "/login"
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
		fmt.Printf("%s\n", reply.Message)
		cfg.Token = reply.Token
		cfg.Authorized = true
		if cfg.Token != "" {
			envMap := map[string]string{"CLI_TOKEN": cfg.Token, "CLI_USER": cfg.User}
			if err := godotenv.Write(envMap, ".env.local"); err != nil {
				log.Println(err)
			}
		}
		return true
	} else {
		fmt.Printf("%s\n", reply.Error)
		cfg.Authorized = false
		return false
	}
}

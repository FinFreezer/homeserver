package functionality

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	//"net/http"
	"os"
	"strings"

	c "github.com/finfreezer/homeserver/internal/config"
)

func changeDirectory(cfg *c.UserConfig) bool {
	type Response struct {
		Message string `json:"reply"`
		Error   string `json:"error,omitempty"`
	}
	responseParams := Response{}
	fullPath := strings.Join(cfg.Args, "%20")
	fullUrl := os.Getenv("DST_SERVER") + os.Getenv("DFLT_PORT") + "/cd/" + fullPath
	resp, err := http.Get(fullUrl)
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()
	newDec := json.NewDecoder(resp.Body)
	newDec.Decode(&responseParams)
	if resp.StatusCode != 200 {
		fmt.Println("Error: " + responseParams.Error)
	} else {
		fmt.Println(responseParams.Message)
	}
	return true
}

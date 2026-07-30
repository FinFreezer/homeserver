package functionality

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	c "github.com/finfreezer/homeserver/internal/config"
)

func remoteUpdate(cfg *c.UserConfig) bool {
	type Response struct {
		Message string `json:"reply"`
		Error   string `json:"error,omitempty"`
	}
	responseParams := Response{}
	fullUrl := os.Getenv("DST_SERVER") + os.Getenv("DFLT_PORT") + "/remoteUpdate"
	req, err := http.NewRequest("POST", fullUrl, nil)
	resp, err := http.DefaultClient.Do(req)
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

package functionality

import (
	"bytes"
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
	type Request struct {
		NewDirectory string `json:"newDir"`
	}
	responseParams := Response{}

	fullPath := strings.Join(cfg.Args, "%20")
	requestBytes, err := json.Marshal(Request{NewDirectory: fullPath})
	if err != nil {
		log.Println(err)
		return false
	}
	requestBuff := bytes.NewBuffer(requestBytes)
	fullUrl := os.Getenv("DST_SERVER") + os.Getenv("DFLT_PORT") + "/cd"
	req, err := http.NewRequest("PUT", fullUrl, requestBuff)
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

package functionality

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	c "github.com/finfreezer/homeserver/internal/config"
)

func listContents(cfg *c.UserConfig) bool {
	fullUrl := ""
	if len(cfg.Args) > 0 && cfg.Args[0] == "-dironly" {
		fullPath := strings.Join(cfg.Args[1:], "%20")
		fullUrl = os.Getenv("DST_SERVER") + os.Getenv("DFLT_PORT") +
			"/listdir/" + fullPath + "?dirOnly=true" + "&recDepth=0"
	} else {
		fullPath := strings.Join(cfg.Args, "%20")
		fullUrl = os.Getenv("DST_SERVER") + os.Getenv("DFLT_PORT") +
			"/listdir/" + fullPath + "?dirOnly=false" + "&recDepth=0"
	}

	log.Printf("Making a GET request to %s\n", fullUrl)
	resp, err := http.Get(fullUrl)
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		newResp := ListDirResponse{}
		decoder := json.NewDecoder(resp.Body)
		if err = decoder.Decode(&newResp); err != nil {
			log.Println(err)
			return false
		}
		fmt.Println(newResp.Message)
		fmt.Println("Directory tree:")
		printContentTree(newResp.Files, 0)
		return true
	} else {
		newResp := ListDirResponse{}
		decoder := json.NewDecoder(resp.Body)
		if err = decoder.Decode(&newResp); err != nil {
			log.Println(err)
			return false
		}
		log.Printf("Something went wrong: %s\n", newResp.Error)
		return false
	}
}

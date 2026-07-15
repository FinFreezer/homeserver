package functionality

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	c "github.com/finfreezer/homeserver/internal/config"
	"github.com/joho/godotenv"
)

type FileInfo struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
}

type ListDirResponse struct {
	Message string     `json:"reply"`
	Files   []FileInfo `json:"directory"`
}

func listcontents(cfg *c.UserConfig) bool {
	fullPath := path.Join(cfg.Args...)
	godotenv.Load()
	fullUrl := os.Getenv("DST_SERVER") + "/listdir/" + fullPath
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
		for _, item := range newResp.Files {
			if item.IsDir {
				fmt.Printf("D - %s\n", item.Name)
			} else {
				fmt.Printf("F - %s\n", item.Name)
			}
		}
		return true
	} else {
		log.Println("Something went wrong.")
		return false
	}

}

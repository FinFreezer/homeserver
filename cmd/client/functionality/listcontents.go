package functionality

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	pg "github.com/finfreezer/homeserver/cmd/server/functionality"
	c "github.com/finfreezer/homeserver/internal/config"
)

type FileInfo struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
}

type ListDirResponse struct {
	Message string      `json:"reply"`
	Files   pg.FileNode `json:"directory"`
	Error   string      `json:"error,omitempty"`
}

func listContents(cfg *c.UserConfig) bool {
	fullPath := path.Join(cfg.Args...)
	fullUrl := os.Getenv("DST_SERVER") + os.Getenv("DFLT_PORT") + "/listdir/" + fullPath
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

func printContentTree(Node pg.FileNode, depth int) {
	if !Node.IsDir {
		for i := 0; i < depth; i++ {
			fmt.Print("\t")
		}
		fmt.Printf("F - %s\n", Node.Name)
	} else {
		for i := 0; i < depth; i++ {
			fmt.Print("\t")
		}
		fmt.Printf("D - %s\n", Node.Name)
		for _, child := range Node.Children {
			printContentTree(child, depth+1)
		}
	}
}

package functionality

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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

func listAllContents(cfg *c.UserConfig) bool {
	fullUrl := ""
	if len(cfg.Args) > 0 && cfg.Args[0] == "-dironly" {
		fullPath := strings.Join(cfg.Args[1:], "%20")
		fullUrl = os.Getenv("DST_SERVER") + os.Getenv("DFLT_PORT") +
			"/listdir/" + fullPath + "?dirOnly=true" + "&recDepth=99"
	} else if len(cfg.Args) > 1 && cfg.Args[0] == "-d" {
		fullPath := strings.Join(cfg.Args[2:], "%20")
		recDepth := cfg.Args[1]
		fullUrl = os.Getenv("DST_SERVER") + os.Getenv("DFLT_PORT") +
			"/listdir/" + fullPath + "?dirOnly=false" + "&recDepth=" + recDepth
	} else {
		fullPath := strings.Join(cfg.Args, "%20")
		fullUrl = os.Getenv("DST_SERVER") + os.Getenv("DFLT_PORT") +
			"/listdir/" + fullPath + "?dirOnly=false" + "&recDepth=99"
	}

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

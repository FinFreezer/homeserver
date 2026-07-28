package functionality

import (
	"log"
	//"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/browser"

	c "github.com/finfreezer/homeserver/internal/config"
)

func streamContent(cfg *c.UserConfig) bool {
	var playType string
	var fullPath string
	if strings.HasPrefix(cfg.Args[0], "-") {
		playType = cfg.Args[0]
		fullPath = strings.Join(cfg.Args[1:], "%20")
	} else {
		playType = ""
		fullPath = strings.Join(cfg.Args, "%20")
	}

	switch playType {
	case "-f":
		fullUrl := os.Getenv("DST_SERVER") + os.Getenv("DFLT_PORT") + "/stream/" + fullPath + "?playlist=false"
		cmd := exec.Command("vlc", fullUrl)
		err := cmd.Run()
		if err != nil {
			log.Println(err)
			return false
		}
		return true
	case "-a":
		fullUrl := os.Getenv("DST_SERVER") + os.Getenv("DFLT_PORT") + "/stream/" + fullPath + "?playlist=true"
		cmd := exec.Command("vlc", fullUrl)
		err := cmd.Run()
		if err != nil {
			log.Println(err)
			return false
		}
		return true
	case "-b":
		fullUrl := os.Getenv("DST_SERVER") + os.Getenv("DFLT_PORT") + "/stream/" + fullPath + "?playlist=false"
		err := browser.OpenURL(fullUrl)
		if err != nil {
			log.Println(err)
			return false
		}
		return true
	default:
		fullUrl := os.Getenv("DST_SERVER") + os.Getenv("DFLT_PORT") + "/stream/" + fullPath + "?playlist=false"
		err := browser.OpenURL(fullUrl)
		if err != nil {
			log.Println(err)
			return false
		}
		return true
	}
}

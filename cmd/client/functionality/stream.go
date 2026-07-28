package functionality

import (
	"log"
	//"net/http"
	"errors"
	"os"
	"os/exec"
	"runtime"
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
		if err := getMediaPlayer(fullUrl); err != nil {
			log.Println(err)
			return false
		}
		return true
	case "-a":
		fullUrl := os.Getenv("DST_SERVER") + os.Getenv("DFLT_PORT") + "/stream/" + fullPath + "?playlist=true"
		if err := getMediaPlayer(fullUrl); err != nil {
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

func getMediaPlayer(pathToStream string) error {
	switch runtime.GOOS {
	case "linux":
		cmd := exec.Command("vlc", pathToStream)
		err := cmd.Run()
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	case "windows":
		cmd := exec.Command("vlc.exe", pathToStream)
		err := cmd.Run()
		if err != nil {
			log.Println(err)
			return err
		}
		return err
	default:
		log.Println("Unsupported environment.")
		return errors.New("Unsupported runtime environment.")
	}
}

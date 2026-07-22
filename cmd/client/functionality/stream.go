package functionality

import (
	"log"
	//"net/http"
	"os"
	"path"

	"github.com/pkg/browser"

	c "github.com/finfreezer/homeserver/internal/config"
)

func streamContent(cfg *c.UserConfig) bool {
	fullPath := path.Join(cfg.Args...)
	fullUrl := os.Getenv("DST_SERVER") + os.Getenv("DFLT_PORT") + "/stream/" + fullPath
	err := browser.OpenURL(fullUrl)
	if err != nil {
		log.Println(err)
		return false
	}
	return true

}

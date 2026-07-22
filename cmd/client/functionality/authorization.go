package functionality

import (
	"errors"

	c "github.com/finfreezer/homeserver/internal/config"
)

type handlerFunc func(cfg *c.UserConfig) bool

func authorizedMiddleware(conf *c.UserConfig, dst handlerFunc, required bool) (bool, error) {
	if conf.Authorized || !required {
		return dst(conf), nil
	} else {
		return false, errors.New("Login required.")
	}
}

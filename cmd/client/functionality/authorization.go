package functionality

import (
	"errors"

	c "github.com/finfreezer/homeserver/internal/config"
)

type handlerFunc func(cfg *c.UserConfig) bool

func authorizedMiddleware(conf *c.UserConfig, dst handlerFunc) (handlerFunc, error) {
	if conf.Authorized {
		return dst, nil
	} else {
		return nil, errors.New("Login required.")
	}
}

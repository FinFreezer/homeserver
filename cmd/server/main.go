package main

import (
	"net/http"

	c "github.com/finfreezer/homeserver/internal/config"
)

func main() {
	apiConfig, _, err := c.OpenDatabase([]string{})
	newMux := http.NewServeMux()
	newMux.HandleFunc("GET /login", a.)
}

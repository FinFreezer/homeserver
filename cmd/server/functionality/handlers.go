package functionality

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/finfreezer/homeserver/internal/auth"
)

type loginParameters struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type response struct {
	Message string `json:"reply"`
}

func (a *ApiConfig) Login(w http.ResponseWriter, r *http.Request) {
	log.Println("Received login request.")
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	params := loginParameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Couldn't decode parameters", err)
		return
	}

	dbUser, err := a.Database.FindUser(context.Background(), params.Name)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find user.", err)
		return
	}

	match, err := auth.CheckPassword(params.Password, dbUser.PasswordHash)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unexpected error.", err)
		return
	}
	if !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect password.", err)
		return
	}

	responseMsg := fmt.Sprintf("Succesfully logged in as %s\n", dbUser.Name)
	respondWithJSON(w, http.StatusOK, response{Message: responseMsg})
}

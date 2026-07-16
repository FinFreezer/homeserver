package functionality

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/finfreezer/homeserver/internal/auth"
)

func (a *ApiConfig) Login(w http.ResponseWriter, r *http.Request) {
	type loginParameters struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	type response struct {
		Message string `json:"reply"`
	}

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
	a.Authorized = true
	respondWithJSON(w, http.StatusOK, response{Message: responseMsg})
}

func (a *ApiConfig) ListContents(w http.ResponseWriter, r *http.Request) {

	type FileInfo struct {
		Name  string `json:"name"`
		IsDir bool   `json:"isDir"`
	}
	type ListDirResponse struct {
		Message string   `json:"reply"`
		Files   FileNode `json:"directory"`
	}

	fullPath := "./assets/" + r.PathValue("path")
	log.Println("Received a request to list contents.")
	fi, err := os.Stat(fullPath)
	if err != nil || !fi.IsDir() {
		log.Println(err)
		respondWithError(w, 400, "Issue finding directory.\n", err)
		return
	}
	resultTree, err := buildFileTree(fullPath)
	if err != nil {
		log.Println(err)
		respondWithError(w, 500, "Issue building directory tree-view.\n", err)
		return
	}
	msg := fmt.Sprintf("Listing files in %s:\n", fullPath)
	respondWithJSON(w, http.StatusOK, ListDirResponse{Message: msg, Files: resultTree})
}

package functionality

import (
	"context"
	"encoding/json"
	"fmt"

	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/finfreezer/homeserver/internal/auth"
	"github.com/finfreezer/homeserver/internal/database"
)

func (a *ApiConfig) Login(w http.ResponseWriter, r *http.Request) {
	type loginParameters struct {
		Name      string `json:"name"`
		Password  string `json:"password"`
		WithToken bool   `json:"withToken"`
		Token     string `json:"token,omitempty"`
	}

	type response struct {
		Message string `json:"reply"`
		Token   string `json:"token,omitempty"`
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
	if params.WithToken {
		userName, err := auth.ValidateJWT(params.Token, a.Secret)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError,
				"Unexpected validation error. Token may be expired, please login.",
				err,
			)
			return
		}
		dbUser, err := a.Database.FindUser(context.Background(), userName)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Couldn't find user.", err)
			return
		}
		responseMsg := fmt.Sprintf("Succesfully logged in as %s\n", dbUser.Name)
		a.Authorized = true
		respondWithJSON(w, http.StatusOK, response{Message: responseMsg})

	} else {
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
		accessToken, err := auth.MakeJWT(dbUser.Name, a.Secret, time.Hour*24*7)
		responseMsg := fmt.Sprintf("Succesfully logged in as %s\n", dbUser.Name)
		params := database.UpdateUserTokenParams{Authtoken: accessToken, Name: dbUser.Name}
		a.Database.UpdateUserToken(context.Background(), params)
		a.Authorized = true
		respondWithJSON(w, http.StatusOK, response{Message: responseMsg, Token: accessToken})
	}
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

func (a *ApiConfig) StreamVideo(w http.ResponseWriter, r *http.Request) {
	fullPath := "./assets/" + r.PathValue("path")
	log.Println("Received a request to stream.")
	fi, err := os.Stat(fullPath)
	if err != nil {
		log.Println(err)
		respondWithError(w, 400, "Couldn't reach file.\n", err)
		return
	}
	if fi.IsDir() {
		respondWithError(w, 400, "Can't stream a directory.\n", err)
		return
	}
	file, err := os.Open(fullPath)
	if err != nil {
		respondWithError(w, 500, "Error opening file.\n", err)
		return
	}
	defer file.Close()
	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Accept-Ranges", "bytes")
	//http.ServeContent(w, r, file.Name(), fi.ModTime(), file) Ready method.
	rangeHeader := r.Header.Get("Range")
	if rangeHeader != "" {
		err, code := readByteRange(file, rangeHeader, w)
		if err != nil && err != io.EOF && code != 0 {
			respondWithError(w, code, "Problem streaming file.", err)
			log.Println(err)
			return
		} else {
			log.Println(err)
			return
		}

	} else {
		w.WriteHeader(200)
		log.Println("Streaming full file.")
		maxBytesPerWrite := 5 << 20
		//maxBytesPerWrite := 5 << 10
		writeData := make([]byte, maxBytesPerWrite)
		var offsetN int64 = 0

		for {
			nFile, err := file.ReadAt(writeData, offsetN)
			if nFile == 0 {
				_, err = w.Write(writeData)
				if err != nil {
					log.Println(err)
					return
				}
				return
			}
			offsetN += int64(nFile)
			//_, err = io.Copy(writeBuffer, file)
			_, err = w.Write(writeData[:nFile])
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

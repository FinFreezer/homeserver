package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	headerValues := headers.Values("Authorization")
	token := make(map[string]string)
	for _, value := range headerValues {
		tokens := strings.Split(value, " ")
		token[tokens[0]] = tokens[1]
	}
	if value, ok := token["ApiKey"]; ok {
		return value, nil
	} else {
		return "", errors.New("No API Key found.")
	}
}

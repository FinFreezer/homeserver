package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	headerValues := headers.Values("Authorization")
	token := make(map[string]string)
	for _, value := range headerValues {
		tokens := strings.Split(value, " ")
		token[tokens[0]] = tokens[1]
	}
	if value, ok := token["Bearer"]; ok {
		return value, nil
	} else {
		return "", errors.New("No token found.")
	}
}

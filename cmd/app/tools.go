package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

func setJsonCookie(w http.ResponseWriter, name string, value interface{}, expiration time.Duration) error {
	cookieValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	escapedValue := url.QueryEscape(string(cookieValue))
	http.SetCookie(w, &http.Cookie{
		Name:    name,
		Value:   escapedValue,
		Path:    "/",
		Expires: time.Now().AddDate(0, 0, 1),
	})
	return nil
}

func getJsonCookie(r *http.Request, name string, value interface{}) error {
	cookie, err := r.Cookie(name)
	if err != nil {
		return err
	}
	decodedValue, err := url.PathUnescape(cookie.Value)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(decodedValue), value)
	if err != nil {
		return err
	}
	return nil
}

func retrieveUserRoom(userID string) (string, bool) {
	for key, userIDSlice := range roomIDDict {
		for _, currUserID := range userIDSlice {
			if currUserID == userID {
				return key, true
			}
		}
	}
	return "", false
}

func removeValueFromSlice(words []string, valueToRemove string) []string {
	var result []string

	for _, word := range words {
		if word != valueToRemove {
			result = append(result, word)
		}
	}

	return result
}

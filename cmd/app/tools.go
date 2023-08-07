package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func delCookieIfExists(w http.ResponseWriter, r *http.Request, name string) {
	_, err := r.Cookie(name)
	if err != nil {
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    name,
		Value:   "",
		Path:    "/",
		Expires: time.Now().AddDate(0, 0, -1),
	})
}

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

func containsInSlice(s []string, item string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == item {
			return true
		}
	}
	return false
}

func getConnectedUsersCount(roomID string) int {
	count := 0
	for _, user := range roomIDDict[roomID] {
		userID, _ := strconv.Atoi(user)
		if userID > 0 {
			count++
		}
	}
	return count
}

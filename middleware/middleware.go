package middleware

import (
	"apex-challenge/db"
	"encoding/json"
	"net/http"
)

func KnownUser(handler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// authenticate with username
		authorizationHeader := request.Header.Get("Authorization")
		_, found := db.GetUser(authorizationHeader)
		if !found {
			unknownUserResponse(writer)
			return
		} else {
			handler.ServeHTTP(writer, request)
		}
	}
}

func ActiveSession(handler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// check for active session
		authorizationHeader := request.Header.Get("Authorization")
		_, _, sFound := db.GetSession(authorizationHeader)
		if sFound {
			handler.ServeHTTP(writer, request)
		} else {
			inactiveSessionResponse(writer)
			return
		}
	}
}

func unknownUserResponse(writer http.ResponseWriter) {
	response := map[string]interface{}{
		"message": "Unknown User",
		"data":    nil,
	}
	writer.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(writer).Encode(response)
}

func inactiveSessionResponse(writer http.ResponseWriter) {
	response := map[string]interface{}{
		"message": "No active game session",
		"data":    nil,
	}
	writer.WriteHeader(http.StatusForbidden)
	json.NewEncoder(writer).Encode(response)
}

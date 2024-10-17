package api

import (
	"encoding/json"
	"log"
)

// Type used for creating a new user in auth
type UserCreation struct {
	password string
	username string
	email    string
}

// Type used for logging in a user in auth
type UserLogin struct {
	password string
	email    string
}

type Message struct {
	Content string `json:"Message"`
}

func ConstructResponse(message string) []byte {
	msg := Message{
		Content: message,
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("Error marshalling to JSON: %s", err)
	}

	return jsonData
}

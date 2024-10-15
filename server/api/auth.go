package api

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"net/http"
	"strconv"
)

func hash(s string) (string, error) {
	h := fnv.New32a()
	_, err := h.Write([]byte("Random Password" + s + "Please Don't be Hacked"))
	if err != nil {
		log.Println("Something went wrong with hashing " + err.Error())
		return "", nil
	}
	return strconv.Itoa(int(h.Sum32())), nil
}
func AuthHandler(jwtSalt []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
		switch r.Method {
		case http.MethodPost:
			var PossibleUser UserCreation
			err := json.NewDecoder(r.Body).Decode(&PossibleUser)
			if err != nil {
				log.Println("Something went wrong with JSON decode " + err.Error())
				http.Error(w, "Something went wrong", http.StatusBadRequest)
			}

			// Hashing password (should already be hashed from user end)

			PossibleUser.password, err = hash(PossibleUser.password)
			if err != nil {
				log.Println("Something went wrong with Hashing password " + err.Error())
				http.Error(w, "Something went wrong", http.StatusBadRequest)
			}
			fmt.Println(PossibleUser.password)
			// 1. Check if username is taken

			// 2. Check if email already is used

			// Return verification token
			token, err := JwtCreation(PossibleUser.username, jwtSalt)
			if err != nil {
				log.Println("Something went wrong with generating token " + err.Error())
				http.Error(w, "Something went wrong", http.StatusBadRequest)
			}
			w.WriteHeader(200)
			_, err = w.Write([]byte(token))
			if err != nil {
				return
			}

		case http.MethodPut:
			var LoginAttempt UserLogin
			err := json.NewDecoder(r.Body).Decode(&LoginAttempt)
			if err != nil {
				log.Println("Something went wrong with JSON decode " + err.Error())
				http.Error(w, "Something went wrong", http.StatusBadRequest)
			}
			_, err = hash(LoginAttempt.password)

			if err != nil {
				log.Println("Something went wrong with Hashing " + err.Error())
				http.Error(w, "Something went wrong", http.StatusBadRequest)
			}
			//Check if password and email match
		default:
			http.Error(w, "This method is not allowed", http.StatusMethodNotAllowed)
		}
	}
}

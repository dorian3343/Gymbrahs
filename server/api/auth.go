package api

import (
	"database/sql"
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
func AuthHandler(jwtSalt []byte, db *sql.DB) http.HandlerFunc {
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
			// Check if username is taken
			usernameExists := false

			err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", PossibleUser.username).Scan(&usernameExists)
			if err != nil {
				panic(err)
			}
			if usernameExists {
				w.WriteHeader(409)
				_, err = w.Write(ConstructResponse("Username is taken"))
				if err != nil {
					return
				}
			} else {
				emailExists := false
				err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", PossibleUser.email).Scan(&emailExists)
				if err != nil {
					panic(err)
				}
				if emailExists {
					w.WriteHeader(409)
					_, err = w.Write(ConstructResponse("Email is takne"))
					if err != nil {

					}
				} else {
					fmt.Println(PossibleUser.username)
					_, err = db.Exec(`INSERT INTO users (username, password, email) VALUES (?, ?, ?)`, PossibleUser.username, PossibleUser.password, PossibleUser.email)
					if err != nil {
						http.Error(w, "Error creating user", http.StatusInternalServerError)
						return
					}

					token, err := JwtCreation(PossibleUser.username, jwtSalt)
					if err != nil {
						log.Println("Something went wrong with generating token " + err.Error())
						http.Error(w, "Something went wrong", http.StatusBadRequest)
					}
					w.WriteHeader(200)
					_, err = w.Write(ConstructResponse(token))
					if err != nil {
						return
					}
				}

			}

			// Check if email  is used

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

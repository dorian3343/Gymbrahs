package api

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

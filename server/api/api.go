package api

import (
	"net/http"
)

// debug endpoint
func GetRoot(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Hewwo :3 Server iz up!! "))
}

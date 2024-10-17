package api

import (
	"net/http"
)

// debug endpoint
func GetRoot(w http.ResponseWriter, _ *http.Request) {
	w.Write(ConstructResponse("Hewwo :3 Server iz up!!"))
}

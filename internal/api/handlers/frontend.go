package handlers

import (
	"net/http"
)

// ServeIndex serves the main frontend page
func ServeIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/static/index.html")
}


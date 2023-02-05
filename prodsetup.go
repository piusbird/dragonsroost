//go:build prod

package main

import "net/http"

func setupDatabase(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not in production bro!", http.StatusForbidden)
	return
}

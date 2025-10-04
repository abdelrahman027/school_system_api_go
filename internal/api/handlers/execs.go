package handlers

import "net/http"

func ExecHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the exec page!"))
}

package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	mw "schoolapi/internal/api/middlewares"
	"strings"
	"time"
)

type user struct {
	Name string `json:"name"`
	Age  string `json:"age"`
	City string `json:"city"`
}

// cmd for creating self-signed certificate
// openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem -config openssl.cnf
func teachersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	id := strings.Split(r.URL.Path, "/")
	fmt.Println("ID:", id)
	fmt.Println("ID reak:", id[2])

	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Welcome to the teachers page! Get method"))
	case http.MethodPost:
		w.Write([]byte("Welcome to the teachers page! Post method"))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
	// w.Write([]byte("Welcome to the teachers page!"))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Hello, World!")
	fmt.Println("Method: ", r.Method)

	w.Write([]byte("Hello, World!"))
}

func execHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Method: ", r.Method)
	w.Write([]byte("Welcome to the exec page!"))
}

func studentsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Method: ", r.Method)
	w.Write([]byte("Welcome to the students page!"))
}
func main() {
	port := ":3000"
	keyFile := "key.pem"
	certFile := "cert.pem"
	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler)

	mux.HandleFunc("/students", studentsHandler)

	mux.HandleFunc("/teachers/", teachersHandler)

	mux.HandleFunc("/exec", execHandler)

	tlsconfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	r1 := mw.NewRateLimiter(5, time.Minute)
	server := &http.Server{
		Addr: port,
		// Handler:   middlewares.SecurityHeaders(mux),
		Handler:   r1.Middleware(mw.Compression(mw.SecurityHeaders(mw.Cors(mw.ResponseTime(mux))))),
		TLSConfig: tlsconfig,
	}

	fmt.Println("Starting server on port", port)
	err := server.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		panic(err)
	}

}

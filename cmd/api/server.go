package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	mw "schoolapi/internal/api/middlewares"
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

	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Welcome to the teachers page! Get method"))
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		fmt.Println("query params:", r.URL.Query())
		fmt.Println("query params:", r.URL.Query().Get("name"))
		fmt.Println("form data:", r.Form)
		w.Write([]byte("Welcome to the teachers page! Post method"))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func execHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the exec page!"))
}

func studentsHandler(w http.ResponseWriter, r *http.Request) {
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

	HppOptions := mw.HppOptions{
		CheckBody:                    true,
		CheckQuery:                   true,
		CheckBodyOnlyForContentTypes: "application/x-www-form-urlencoded",
		Whitelist:                    []string{"allowed"},
	}
	// secureMux := mw.Hpp(HppOptions)(r1.Middleware(mw.Compression(mw.SecurityHeaders(mw.Cors(mw.ResponseTime(mux))))))
	secureMux := applyMiddlewares(mux,
		mw.Hpp(HppOptions),
		mw.Compression,
		mw.SecurityHeaders,
		mw.ResponseTime,
		r1.Middleware,
		mw.Cors,
	)
	server := &http.Server{
		Addr:      port,
		Handler:   secureMux,
		TLSConfig: tlsconfig,
	}

	fmt.Println("Starting server on port", port)
	err := server.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		panic(err)
	}

}

type Middleware func(http.Handler) http.Handler

func applyMiddlewares(handler http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}

package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	mw "schoolapi/internal/api/middlewares"
	"strconv"
	"strings"
)

type Teacher struct {
	ID        int
	FirstName string
	LastName  string
	Class     string
	Subject   string
}

var (
	teachers = make(map[int]Teacher)
	// Mutex = &sync.Mutex{}
	nextID = 1
)

func init() {
	teachers[nextID] = Teacher{ID: nextID, FirstName: "John", LastName: "Doe", Class: "10A", Subject: "Math"}
	nextID++
	teachers[nextID] = Teacher{ID: nextID, FirstName: "Jane", LastName: "Smith", Class: "10B", Subject: "Science"}
	nextID++
	teachers[nextID] = Teacher{ID: nextID, FirstName: "Emily", LastName: "Doe", Class: "10C", Subject: "History"}
	nextID++
	teachers[nextID] = Teacher{ID: nextID, FirstName: "Michael", LastName: "Brown", Class: "10D", Subject: "English"}
	nextID++
	teachers[nextID] = Teacher{ID: nextID, FirstName: "Sarah", LastName: "Davis", Class: "10E", Subject: "Art"}
	nextID++
	teachers[nextID] = Teacher{ID: nextID, FirstName: "David", LastName: "Wilson", Class: "10F", Subject: "Physical Education"}
}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	IDstr := strings.TrimSuffix(path, "/")
	fmt.Println("IDstr:", IDstr)
	if IDstr == "" {
		firstName := r.URL.Query().Get("first_name")
		lastName := r.URL.Query().Get("last_name")

		teachersList := make([]Teacher, 0, len(teachers))
		for _, teacher := range teachers {
			if (firstName == "" || teacher.FirstName == firstName) && (lastName == "" || teacher.LastName == lastName) {
				teachersList = append(teachersList, teacher)
			}
		}
		response := struct {
			Status string    `json:"status"`
			Count  int       `json:"count"`
			Data   []Teacher `json:"data"`
		}{
			Status: "success",
			Count:  len(teachers),
			Data:   teachersList,
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}

	}
	numID, err := strconv.Atoi(IDstr)
	if err != nil {
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}
	teacher, exists := teachers[numID]
	if !exists {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(teacher)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// cmd for creating self-signed certificate
// openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem -config openssl.cnf
func teachersHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		getTeachersHandler(w, r)
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
	// r1 := mw.NewRateLimiter(5, time.Minute)

	// HppOptions := mw.HppOptions{
	// 	CheckBody:                    true,
	// 	CheckQuery:                   true,
	// 	CheckBodyOnlyForContentTypes: "application/x-www-form-urlencoded",
	// 	Whitelist:                    []string{"allowed"},
	// }
	// secureMux := mw.Hpp(HppOptions)(r1.Middleware(mw.Compression(mw.SecurityHeaders(mw.Cors(mw.ResponseTime(mux))))))
	// secureMux := applyMiddlewares(mux,
	// 	mw.Hpp(HppOptions),
	// 	mw.Compression,
	// 	mw.SecurityHeaders,
	// 	mw.ResponseTime,
	// 	r1.Middleware,
	// 	mw.Cors,
	// )

	secureMux := applyMiddlewares(mux,
		mw.SecurityHeaders,
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

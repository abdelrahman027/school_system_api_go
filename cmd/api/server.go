package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	mw "schoolapi/internal/api/middlewares"
	"schoolapi/internal/api/router"
	"schoolapi/internal/repository/sqlconnect"
	"schoolapi/pkg/utils"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	_, err := sqlconnect.ConnectDB(os.Getenv("DB_NAME"))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	port := os.Getenv("API_PORT")
	// cmd for creating self-signed certificate
	// openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem -config openssl.cnf
	keyFile := "key.pem"
	certFile := "cert.pem"

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
	router := router.Router()
	secureMux := utils.ApplyMiddlewares(router,
		mw.SecurityHeaders,
	)
	server := &http.Server{
		Addr:      port,
		Handler:   secureMux,
		TLSConfig: tlsconfig,
	}

	fmt.Println("Starting server on port", port)
	err = server.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		panic(err)
	}

}

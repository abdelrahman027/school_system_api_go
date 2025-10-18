package main

import (
	"crypto/tls"
	"embed"
	"fmt"
	"net/http"
	"os"
	mw "schoolapi/internal/api/middlewares"
	"schoolapi/internal/api/router"
	"schoolapi/pkg/utils"
	"time"

	"github.com/joho/godotenv"
)

var envFile embed.FS

func main() {
	godotenv.Load()

	// _, err := sqlconnect.ConnectDB(os.Getenv("DB_NAME"))
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	port := os.Getenv("API_PORT")
	// cmd for creating self-signed certificate
	// openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem -config openssl.cnf
	keyFile := "key.pem"
	certFile := "cert.pem"

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
	router := router.MainRouter()
	jwtMiddleware := mw.ExludePaths(mw.JWTMiddleware, "/execs/login", "/execs/forgetpassword", "/execs/resetpassword/reset")
	// secureMux := mw.Hpp(HppOptions)(r1.Middleware(mw.Compression(mw.SecurityHeaders(mw.Cors(mw.ResponseTime(mux))))))
	secureMux := utils.ApplyMiddlewares(router,
		mw.Compression,
		mw.SecurityHeaders,
		mw.Hpp(HppOptions),
		mw.XSSMiddleware,
		jwtMiddleware,
		mw.ResponseTime,
		r1.Middleware,
		mw.Cors,
	)
	// secureMux := utils.ApplyMiddlewares(router,
	// 	mw.SecurityHeaders,
	// 	jwtMiddleware,

	// )
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

package utils

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func CheckHttpError(err error, w http.ResponseWriter, errorMsg string, httpStatus int) {
	if err != nil {
		log.Println("Error: ", err)
		http.Error(w, errorMsg, httpStatus)
	}

}

func ErrorHandler(err error, message string) error {
	errorLogger := log.New(os.Stderr, "Error: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger.Println(message, err)
	return fmt.Errorf(message)
}

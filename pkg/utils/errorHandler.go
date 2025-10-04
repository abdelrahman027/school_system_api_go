package utils

import (
	"log"
	"net/http"
)

func CheckHttpError(err error, w http.ResponseWriter, errorMsg string, httpStatus int) {
	if err != nil {
		log.Println("Error: ", err)
		http.Error(w, errorMsg, httpStatus)
	}

}

package router

import (
	"net/http"
	"schoolapi/internal/api/handlers"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /teachers/", handlers.GetTeacherHandler)
	mux.HandleFunc("POST /teachers/", handlers.CreateTeacherHandler)
	mux.HandleFunc("PATCH /teachers/", handlers.PatchTeacherHandler)
	mux.HandleFunc("DELETE /teachers/", handlers.DeleteTeacherHandler)

	mux.HandleFunc("GET /teachers/{id}", handlers.GetOneTeacherHandler)
	mux.HandleFunc("PUT /teachers/{id}", handlers.UpdateTeacherHandler)
	mux.HandleFunc("PATCH /teachers/{id}", handlers.PatchOneTeacherHandler)
	mux.HandleFunc("DELETE /teachers/{id}", handlers.DeleteOneTeacherHandler)

	mux.HandleFunc("/", handlers.RootHandler)

	mux.HandleFunc("/students/", handlers.StudentsHandler)

	mux.HandleFunc("/exec/", handlers.ExecHandler)
	return mux

}

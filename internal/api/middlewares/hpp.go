package middlewares

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
)

// define options for HPP middleware
type HppOptions struct {
	// Add fields if necessary
	CheckQuery                   bool
	CheckBody                    bool
	CheckBodyOnlyForContentTypes string
	Whitelist                    []string
}

// middleware that take a parameter
func Hpp(options HppOptions) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if options.CheckBody && r.Method == http.MethodPost && isCorrectContentType(r, options.CheckBodyOnlyForContentTypes) {
				//filter body parameters
				filterparams(r, options.Whitelist)
			}
			if options.CheckQuery && r.URL.Query() != nil {
				//filter query parameters
				filterQuerys(r, options.Whitelist)
			}
			// Implement HPP protection logic here
			// For example, you can check for duplicate parameters in the query string or body
			// and remove them if they are not in the whitelist
			next.ServeHTTP(w, r)
		})
	}
}

// to define if content type is correct
func isCorrectContentType(r *http.Request, contentType string) bool {
	return strings.Contains(r.Header.Get("Content-Type"), contentType)
}

// main filtering function to remove duplicate parameters
func filterparams(r *http.Request, whitelist []string) {
	// Implement parameter filtering logic here
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		return
	}
	for k, v := range r.Form {
		if len(v) > 1 {
			//remove duplicate
			r.Form.Set(k, v[0])
		}
		if !isInWhitelist(k, whitelist) {
			//remove parameter
			delete(r.Form, k)
		}
	}

}

// check if parameter is in whitelist
func isInWhitelist(param string, whitelist []string) bool {
	return slices.Contains(whitelist, param)
}

// main Query filtering function to remove duplicate parameters
func filterQuerys(r *http.Request, whitelist []string) {
	queries := r.URL.Query()

	for k, v := range queries {
		if len(v) > 0 {
			//remove duplicate
			queries.Set(k, v[0])
		}
		if !isInWhitelist(k, whitelist) {
			//remove parameter
			queries.Del(k)
		}
	}
	r.URL.RawQuery = queries.Encode()
}

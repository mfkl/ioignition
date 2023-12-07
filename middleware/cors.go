package middleware

import (
	"net/http"

	"github.com/go-chi/cors"
)

var options = cors.Options{
	AllowedOrigins: []string{"*"},
	AllowedMethods: []string{"GET", "PUT", "POST", "DELETE", "HEAD", "OPTION"},
	AllowedHeaders: []string{
		"User-Agent",
		"Content-Type",
		"Accept",
		"Accept-Encoding",
		"Accept-Language",
		"Cache-Control",
		"Connection",
		"DNT",
		"Host",
		"Origin",
		"Pragma",
		"Referer",
	},
	ExposedHeaders: []string{"Link"},
	MaxAge:         300, // Maximum value not ignored by any of major browsers
}

func Cors() func(next http.Handler) http.Handler {
	return cors.Handler(options)
}

func ExternalApiCors() func(next http.Handler) http.Handler {
	externalOptions := options
	externalOptions.AllowedMethods = []string{"POST", "PUT", "OPTION"}

	return cors.Handler(externalOptions)
}

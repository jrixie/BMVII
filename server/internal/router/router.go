package router

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"boilermakevii/api/internal/alert"
	"boilermakevii/api/internal/test"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// NewRouter constructs a new router for the API to route requests.
func NewRouter() (router *mux.Router) {
	router = mux.NewRouter().StrictSlash(true)
	router.Use(apiMiddleware)

	for _, service := range apiSpec {
		for _, route := range service.Routes {
			handler := route.Handler

			path := fmt.Sprintf(
				"/api/%s/%s/",
				strings.Trim(service.Pattern, "/"),
				strings.Trim(route.Pattern, "/"),
			)

			router.
				Name(route.Name).
				Methods("OPTIONS", route.Method).
				Path(path).
				Handler(handler)

		}
	}

	return router
}

// apiMiddleware defines a common middleware used across all endpoints
func apiMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		next.ServeHTTP(w, r)
		log.Info(fmt.Sprintf("%s\t%s\tin %v", r.Method, r.URL.Path, time.Since(start)))
	})
}

// Service defines a structure for API services.
type Service struct {
	Pattern string
	Routes  []Route
}

// Route defines a structure for endpoints in a Service.
type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler http.HandlerFunc
}

// apiSpec is the specification for the API. Add new services and their routes here.
var apiSpec = []Service{
	{
		"test",
		[]Route{
			{
				"HelloWorld",
				"GET",
				"/",
				test.HelloWorld,
			},
		},
	},
	{
		"alert",
		[]Route{
			{
				"CreateTrigger",
				"POST",
				"/create",
				alert.CreateTrigger,
			},
		},
	},
}

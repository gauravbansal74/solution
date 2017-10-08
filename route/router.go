package route

import (
	"github.com/gauravbansal74/solution/handler"
	"github.com/gorilla/mux"
	"net/http"
)

func newRouter() *mux.Router {
	router := mux.NewRouter()
	// Server Info end Point to check server information
	router.Methods("GET").Path("/info").Name("Info").HandlerFunc(handler.Info)
	router.Methods("GET").Path("/route/{token}").Name("RouteGet").HandlerFunc(handler.Get)
	router.Methods("POST").Path("/route").Name("RouteGet").HandlerFunc(handler.Post)
	return router
}

// Load the HTTP routes and validate user token
func LoadHandler() http.Handler {
	return &WithCORS{newRouter()}
}

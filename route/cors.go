package route

import (
	"github.com/gorilla/mux"
	"net/http"
)

type WithCORS struct {
	r *mux.Router
}

// TODO: This needs to be fixed with per-handler CORS treatments. Only a placeholder
func (s *WithCORS) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		res.Header().Set("Access-Control-Allow-Origin", origin)
		res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		res.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Dragon-Law-Username, X-Dragon-Law-Dragonball, X-Dragon-Law-API-Version")
	}

	// Stop here for a Preflighted OPTIONS request.
	if req.Method == "OPTIONS" {
		return
	}
	// Lets Gorilla work
	s.r.ServeHTTP(res, req)
}

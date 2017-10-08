package server

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

// Server stores the hostname and port number
type Server struct {
	ListenHost string `json:"ListenHost"` // Server name
	ListenPort int32  `json:"ListenPort"` // HTTP port
}

// startHTTP starts the HTTP listener
func StartServer(handlers http.Handler, s Server) {

	log.Info("Server Started. Listening on " + httpAddress(s))
	// Start the HTTP listener
	log.Fatal(http.ListenAndServe(httpAddress(s), handlers))
}

// httpAddress returns the HTTP address
func httpAddress(s Server) string {
	return s.ListenHost + ":" + fmt.Sprintf("%d", s.ListenPort)
}

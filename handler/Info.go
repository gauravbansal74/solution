package handler

import (
	// "fmt"
	"github.com/gauravbansal74/solution/config"
	"net/http"
)

var (
	conf config.Config
)

func Info(w http.ResponseWriter, r *http.Request) {
	conf = config.ReadConfig()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(conf.SystemInfo))
	return
}

package response

import (
	"encoding/json"
	"net/http"
)

// Core Response
type Token struct {
	ID string `json:"token"`
}

type Success struct {
	Status        interface{} `json:"status"`
	Path          interface{} `json:"path"`
	TotalDistance interface{} `json:"total_distance"`
	TotalTime     interface{} `json:"total_time"`
}

// Core Response
type Other struct {
	Message interface{} `json:"status"`
}

// Change Response
type Error struct {
	Message interface{} `json:"error"`
	Status  interface{} `json:"status"`
}

// SendError calls Send by without a count or results
func SendError(w http.ResponseWriter, status http.ConnState, message interface{}) {
	i := Error{}
	i.Message = message
	i.Status = "failure"

	js, err := json.Marshal(i)
	if err != nil {
		http.Error(w, "JSON Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(status))
	w.Write(js)
}

// SendToken calls Send by without a count or results
func SendToken(w http.ResponseWriter, tok string) {
	i := Token{}
	i.ID = tok

	js, err := json.Marshal(i)
	if err != nil {
		http.Error(w, "JSON Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(http.StatusOK))
	w.Write(js)
}

// SendOther calls Send by without a count or results
func SendOther(w http.ResponseWriter, status http.ConnState, message interface{}) {
	i := Other{}
	i.Message = message

	js, err := json.Marshal(i)
	if err != nil {
		http.Error(w, "JSON Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(status))
	w.Write(js)
}

// SendSuccess calls Send by without a count or results
func SendSuccess(w http.ResponseWriter, Status interface{}, Path interface{}, TotalDistance interface{}, TotalTime interface{}) {
	i := Success{}
	i.Status = Status
	i.Path = Path
	i.TotalDistance = TotalDistance
	i.TotalTime = TotalTime

	js, err := json.Marshal(i)
	if err != nil {
		http.Error(w, "JSON Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(http.StatusOK))
	w.Write(js)
}

// Send writes struct to the writer using a format
func Send(w http.ResponseWriter, status http.ConnState, results interface{}) {

	js, err := json.Marshal(results)
	if err != nil {
		http.Error(w, "JSON Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(status))
	w.Write(js)
}

// SendJSON writes a struct to the writer
func SendJSON(w http.ResponseWriter, i interface{}) {
	js, err := json.Marshal(i)
	if err != nil {
		http.Error(w, "JSON Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

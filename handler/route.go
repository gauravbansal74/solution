package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	// sdistance "github.com/gauravbansal74/solution/distance"
	"github.com/gauravbansal74/solution/model/distance"
	"github.com/gauravbansal74/solution/shared/queue"
	"github.com/gauravbansal74/solution/shared/response"
	"github.com/gorilla/mux"
)

type Message struct {
	Status string `json:"status"`
}

func Post(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "Error while reading Body data")
		return
	}
	if string(body) == "" {
		response.SendError(w, http.StatusBadRequest, "Body data can't be null or empty")
		return
	}

	bData := [][]string{}
	err = json.Unmarshal(body, &bData)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "Error while parse JSON Body data into struct values")
		return
	}

	if len(bData) > 0 {
		// Push to redis
		distanceObject, err := distance.New()
		if err != nil {
			response.SendError(w, http.StatusBadRequest, err.Error())
			return
		}
		distanceObject.Path = bData
		queueConnection := queue.ReadQueue()
		bytesse, err := json.Marshal(distanceObject)
		if err != nil {
			response.SendError(w, http.StatusBadRequest, err.Error())
			return
		}
		err = queueConnection.PutPayload(string(bytesse))
		if err != nil {
			response.SendError(w, http.StatusBadRequest, err.Error())
			return
		} else {
			redisapi := queue.ReadRedisClient()
			err = redisapi.Client.Set(distanceObject.ID, bytesse, 0).Err()
			if err != nil {
				response.SendError(w, http.StatusBadRequest, "Error while saving message in redis")
				return
			}
			response.SendToken(w, distanceObject.ID)
			return
		}
	} else {
		response.SendError(w, http.StatusBadRequest, "Invalid data")
		return
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	redisapi := queue.ReadRedisClient()
	redisData, err := redisapi.Client.Get(token).Result()
	if err != nil {
		if redisData == "" {
			response.SendError(w, http.StatusBadRequest, "Token Invaild")
			return
		} else {
			response.SendError(w, http.StatusBadRequest, err.Error())
			return
		}

	}
	responseData := distance.Entity{}
	err = json.Unmarshal([]byte(redisData), &responseData)
	if err != nil {
		fmt.Println(responseData)
		response.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println(responseData)
	switch responseData.Status {
	case 0:
		{
			fmt.Println(responseData)
			response.SendError(w, http.StatusBadRequest, responseData.Message)
			return
		}
	case 1:
		{
			response.SendSuccess(w, responseData.Message, responseData.Path, responseData.TotalDistance, responseData.TotalTime)
			return
		}
	case 2:
		{
			response.SendOther(w, http.StatusOK, responseData.Message)
			return
		}
	}
	response.Send(w, http.StatusOK, redisData)
	return
}

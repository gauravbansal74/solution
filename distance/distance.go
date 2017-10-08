package distance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gauravbansal74/solution/config"
	mdistance "github.com/gauravbansal74/solution/model/distance"
	"github.com/gauravbansal74/solution/shared/logger"
)

const (
	URL   = "https://maps.googleapis.com/maps/api/distancematrix/json"
	Units = "imperial"
)

type DistanceResult struct {
	DestinationAddresses []string `json:"destination_addresses"`
	OriginAddresses      []string `json:"origin_addresses"`
	Rows                 []struct {
		Elements []struct {
			Distance struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"distance"`
			Duration struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"duration"`
			Status string `json:"status"`
		} `json:"elements"`
	} `json:"rows"`
	Status string `json:"status"`
}

func GetMapDistance(origin []string, destination []string) (DistanceResult, error) {
	dResult := DistanceResult{}
	// Prepare Request URI with query parameters
	apiURI, err := url.Parse(URL)
	if err != nil {
		return dResult, err
	}
	conf := config.LoadConfig()
	query := apiURI.Query()
	query.Add("units", Units)
	query.Add("key", conf.GoogleApiKey)
	query.Add("origins", strings.Join(origin, ","))
	query.Add("destinations", strings.Join(destination, ","))
	apiURI.RawQuery = query.Encode()

	// Make http Get request to get result using google distance API.
	resp, err := http.Get(apiURI.String())
	if err != nil {
		return dResult, err
	}
	defer resp.Body.Close()
	// Read Response Body bytes
	body, err := ioutil.ReadAll(resp.Body)
	// unmarshal JSON into an struct
	err = json.Unmarshal(body, &dResult)
	if err != nil {
		return dResult, err
	}
	return dResult, nil
}

func UniquePathways(length int) []string {
	allLocations := map[int]int{}
	uniquePath := make([]string, length-1)
	// Create total destination location Array
	for i := 1; i < length; i++ {
		allLocations[i] = i
	}
	for i := 1; i < length; i++ {
		for j := 1; j < length-1; j++ {
			temp := allLocations[j]
			allLocations[j] = allLocations[j+1]
			allLocations[j+1] = temp
			currentPath := "0"
			for k := 1; k < length; k++ {
				currentPath = currentPath + strconv.Itoa(allLocations[k])
			}
			uniquePath[i-1] = currentPath
		}
	}
	return uniquePath
}

func CheckForShortestDistance(payload string) (*mdistance.Entity, error) {
	var err error
	model, err := mdistance.Parse(payload)
	if err != nil {
		logger.Error("server", err, "Error while parsing redis queue message", logger.Fields{
			"Id": model.ID,
		})
		return model, err
	} else {
		length := len(model.Path)
		uniquePaths := UniquePathways(length)
		if len(uniquePaths) > 0 {
			values := map[string]int{}
			timeValues := map[string]int{}
			for i := 0; i < len(model.Path); i++ {
				startFrom := model.Path[i]
				if len(startFrom) == 2 {
					for j := 1; j < len(model.Path); j++ {
						if i != j { // If i == j, it means distance is zero between same locations.
							endTo := model.Path[j]
							if len(endTo) == 2 {
								if startFrom[0] == endTo[0] && startFrom[1] == endTo[1] {
									// If i !=j but still long and lat is same for locations then distance would be zero.
									values[strconv.Itoa(i)+"_"+strconv.Itoa(j)] = 0
									timeValues[strconv.Itoa(i)+"_"+strconv.Itoa(j)] = 0
								} else {
									mapDistance, err := GetMapDistance(startFrom, endTo)
									if err != nil {
										return model, err
									} else {
										if len(mapDistance.Rows) > 0 {
											values[strconv.Itoa(i)+"_"+strconv.Itoa(j)] = mapDistance.Rows[0].Elements[0].Distance.Value
											timeValues[strconv.Itoa(i)+"_"+strconv.Itoa(j)] = mapDistance.Rows[0].Elements[0].Duration.Value
										} else {
											logger.Error("server", fmt.Errorf("Error while fetching distance using Google maps"), "Error while fetching distance using Google maps", logger.Fields{
												"Id": model.ID,
											})
											return model, fmt.Errorf("Error while fetching distance using Google maps")
										}
									}
								}
							} else {
								logger.Error("server", fmt.Errorf("Long Lat is not proper"), "Long Lat is not proper", logger.Fields{
									"Id": model.ID,
								})
								return model, fmt.Errorf("Long Lat is not proper")
							}
						} else {
							// If i and j are equals then distance would be zero always
							values[strconv.Itoa(i)+"_"+strconv.Itoa(j)] = 0
							timeValues[strconv.Itoa(i)+"_"+strconv.Itoa(j)] = 0
						}
					}
				} else {
					logger.Error("server", fmt.Errorf("Long Lat is not proper"), "Long Lat is not proper", logger.Fields{
						"Id": model.ID,
					})
					return model, fmt.Errorf("Long Lat is not proper")
				}
			}
			values["0_0"] = 0
			shorestPath := CalculatePathwaysDistance(values, uniquePaths)
			outputData := PrepareOutput(shorestPath, values, timeValues, *model)
			return &outputData, nil
		} else {
			logger.Error("server", fmt.Errorf("Unique Path length is Zero"), "Unique pathways can't be zero", logger.Fields{
				"Id": model.ID,
			})
			return model, fmt.Errorf("Unique pathways can't be zero")
		}
	}
	return model, nil
}

func PrepareOutput(shorestPath string, allDistanceValues map[string]int, allTimeValues map[string]int, model mdistance.Entity) mdistance.Entity {
	previeousPaths := make([][]string, len(model.Path))
	for i := 0; i < len(model.Path); i++ {
		previeousPaths[i] = model.Path[i]
	}
	for i := 0; i <= len(shorestPath)-1; i++ {

		if i < len(shorestPath)-1 {
			model.TotalDistance = model.TotalDistance + allDistanceValues[string(shorestPath[i])+"_"+string(shorestPath[i+1])]
			model.TotalTime = model.TotalTime + allTimeValues[string(shorestPath[i])+"_"+string(shorestPath[i+1])]
		}
		index, _ := strconv.ParseInt(string(shorestPath[i]), 10, 0)
		model.Path[i] = previeousPaths[index]
	}
	model.Status = 1
	model.Message = "success"
	return model
}

func CalculatePathwaysDistance(allDistances map[string]int, uniquePathways []string) string {
	allDistanceValues := map[string]int{}
	for i := 0; i < len(uniquePathways); i++ {
		uniquePath := uniquePathways[i]
		allDistanceValues[uniquePath] = 0
		for j := 0; j < len(uniquePath)-1; j++ {
			allDistanceValues[uniquePath] = allDistanceValues[uniquePath] + allDistances[string(uniquePath[j])+"_"+string(uniquePath[j+1])]
		}
	}
	shorestPath := GetShorestPath(allDistanceValues)
	return shorestPath
}

func GetShorestPath(allDistanceValues map[string]int) string {
	allValues := make([]int, len(allDistanceValues))
	var i = 0
	for _, v := range allDistanceValues {
		allValues[i] = v
		i++
	}
	for k, v := range allDistanceValues {
		if v == allValues[0] && k != "" {
			return k
		}
	}
	return ""
}

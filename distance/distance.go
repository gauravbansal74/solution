package distance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gauravbansal74/solution/config"
	mdistance "github.com/gauravbansal74/solution/model/distance"
	"github.com/gauravbansal74/solution/shared/logger"
)

const (
	URL                = "https://maps.googleapis.com/maps/api/distancematrix/json"
	Units              = "imperial"
	InvalidLongLat     = "Long/Lat are not valid for locations"
	DistanceAPIError   = "Google Maps Distance API error"
	UniquePathwaysZero = "Unique pathways can't be zero"
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

// Get Distance between locations using Google Map Distance API
func getMapDistance(origin []string, destination []string) (DistanceResult, error) {
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

// Find Unique Pathways for multple dropzone locations
func uniquePathways(length int) []string {
	allLocations := map[int]int{}
	var x = 0
	uniquePath := make([]string, int(math.Exp2(float64(length-1))))
	// Create total destination location Array
	for i := 1; i < length; i++ {
		allLocations[i] = i
	}
	for i := 1; i <= length; i++ {
		for j := 1; j < length-1; j++ {
			temp := allLocations[j]
			allLocations[j] = allLocations[j+1]
			allLocations[j+1] = temp
			currentPath := "0"
			for k := 1; k < length; k++ {
				currentPath = currentPath + strconv.Itoa(allLocations[k])
			}
			uniquePath[x] = currentPath
			x++
		}
	}
	return uniquePath
}

// Check  string value is Number or Not -if Number return TRUE else FALSE
func isNumber(input string) bool {
	if _, err := strconv.ParseFloat(input, 64); err == nil {
		return true
	} else {
		return false
	}
}

//
func CheckForShortestDistance(payload string) (*mdistance.Entity, error) {
	var err error
	model, err := mdistance.Parse(payload)
	if err != nil {
		logger.Error("server", err, "Error while parsing redis queue message", logger.Fields{
			"Id": model.ID,
		})
		model.Status = 0
		model.Message = err.Error()
		return model, err
	} else {
		// fmt.Println("initial Model", model)
		length := len(model.Path)
		uniquePaths := uniquePathways(length)
		// fmt.Println("AllUniquePath", uniquePaths)
		if len(uniquePaths) > 0 {
			values := map[string]int{}
			timeValues := map[string]int{}
			for i := 0; i < len(model.Path); i++ {
				startFrom := model.Path[i]
				// fmt.Println("startFrom", startFrom)
				if len(startFrom) == 2 && (isNumber(startFrom[0]) && isNumber(startFrom[1])) {
					for j := 1; j < len(model.Path); j++ {
						if i != j { // If i == j, it means distance is zero between same locations.
							endTo := model.Path[j]
							// fmt.Println("endTo", endTo)
							if len(endTo) == 2 {
								if startFrom[0] == endTo[0] && startFrom[1] == endTo[1] {
									// If i !=j but still long and lat is same for locations then distance would be zero.
									values[strconv.Itoa(i)+"_"+strconv.Itoa(j)] = 0
									timeValues[strconv.Itoa(i)+"_"+strconv.Itoa(j)] = 0
								} else {
									mapDistance, err := getMapDistance(startFrom, endTo)
									if err != nil {
										model.Status = 0
										model.Message = err.Error()
										return model, err
									} else {
										if len(mapDistance.Rows) > 0 {
											values[strconv.Itoa(i)+"_"+strconv.Itoa(j)] = mapDistance.Rows[0].Elements[0].Distance.Value
											timeValues[strconv.Itoa(i)+"_"+strconv.Itoa(j)] = mapDistance.Rows[0].Elements[0].Duration.Value
										} else {
											logger.Error("server", fmt.Errorf(DistanceAPIError), DistanceAPIError, logger.Fields{
												"Id": model.ID,
											})
											model.Status = 0
											model.Message = DistanceAPIError
											return model, fmt.Errorf(DistanceAPIError)
										}
									}
								}
							} else {
								logger.Error("server", fmt.Errorf(InvalidLongLat), InvalidLongLat, logger.Fields{
									"Id": model.ID,
								})
								model.Status = 0
								model.Message = InvalidLongLat
								return model, fmt.Errorf(InvalidLongLat)
							}
						} else {
							// If i and j are equals then distance would be zero always
							values[strconv.Itoa(i)+"_"+strconv.Itoa(j)] = 0
							timeValues[strconv.Itoa(i)+"_"+strconv.Itoa(j)] = 0
						}
					}
				} else {
					logger.Error("server", fmt.Errorf(InvalidLongLat), InvalidLongLat, logger.Fields{
						"Id": model.ID,
					})
					model.Status = 0
					model.Message = InvalidLongLat
					return model, fmt.Errorf(InvalidLongLat)
				}
			}
			values["0_0"] = 0
			// fmt.Println("Values", values)
			// fmt.Println("TimeValues", timeValues)
			// fmt.Println("Models", model)
			shorestPath := calculatePathwaysDistance(values, uniquePaths)
			// fmt.Println("shorestPath", shorestPath)
			outputData := prepareOutput(shorestPath, values, timeValues, *model)
			// fmt.Println("outputData", outputData)
			return &outputData, nil
		} else {
			logger.Error("server", fmt.Errorf(UniquePathwaysZero), UniquePathwaysZero, logger.Fields{
				"Id": model.ID,
			})
			return model, fmt.Errorf(UniquePathwaysZero)
		}
	}
	return model, nil
}

func prepareOutput(shorestPath string, allDistanceValues map[string]int, allTimeValues map[string]int, model mdistance.Entity) mdistance.Entity {
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

func calculatePathwaysDistance(allDistances map[string]int, uniquePathways []string) string {
	allDistanceValues := map[string]int{}
	for i := 0; i < len(uniquePathways); i++ {
		uniquePath := uniquePathways[i]
		allDistanceValues[uniquePath] = 0
		for j := 0; j < len(uniquePath)-1; j++ {
			allDistanceValues[uniquePath] = allDistanceValues[uniquePath] + allDistances[string(uniquePath[j])+"_"+string(uniquePath[j+1])]
		}
	}
	shorestPath := getShorestPath(allDistanceValues)
	return shorestPath
}

func getShorestPath(allDistanceValues map[string]int) string {
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

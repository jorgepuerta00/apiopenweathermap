package common

import (
	"encoding/json"
	"io/ioutil"
)

// APIConfigData struct of OpenWeatherMapApiKey
type APIConfigData struct {
	OpenWeatherMapAPIKey string `json:"OpenWeatherMapApiKey"`
}

// LoadAPIConfig return configuration idapp
func LoadAPIConfig() (APIConfigData, error) {
	filename := "app.config"
	bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		return APIConfigData{}, err
	}

	var c APIConfigData
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return APIConfigData{}, err
	}

	return c, nil
}

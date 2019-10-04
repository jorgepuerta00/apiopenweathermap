package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"weatherProject2/common"
	"weatherProject2/models"
)

// WeatherAPI method exposed
func WeatherAPI() {
	http.HandleFunc("/weather/", handlerGetWeatherInfo)
	http.ListenAndServe(":8080", nil)
}

func handlerGetWeatherInfo(w http.ResponseWriter, r *http.Request) {
	city := strings.SplitN(r.URL.Path, "/", 3)[2]
	data, err := GetWeatherInfo(city)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(data))
}

// GetWeatherInfo Save each request in a database with a timestamp, Load the information from the persistence layer when it’s available. Store new values when the timestamp difference is higher than 300 seconds
func GetWeatherInfo(city string) (string, error) {
	var weather string

	data, err := common.SelectDataWeather(city)
	if err != nil {
		return err.Error(), err
	}

	duration := common.GetTimeElapse(data.Timestamp)

	if int(duration.Seconds()) > 300 {
		weather, err = GetAndSaveWeatherData(city)
		if err != nil {
			return err.Error(), err
		}
	} else {
		weather = models.ReadWeatherDataStruct(data)
	}

	return weather, nil
}

// GetAndSaveWeatherData get and save data from api openweathermap
func GetAndSaveWeatherData(city string) (string, error) {

	data, err := GetDataOpenWeatherMap(city)
	if err != nil {
		return err.Error(), err
	}

	if data.Message == "" {
		err := common.DeleteDataWeather(data)
		if err != nil {
			return err.Error(), err
		}
		err = common.InsertDataWeather(data)
		if err != nil {
			return err.Error(), err
		}
	} else {
		return data.Message, fmt.Errorf(data.Message)
	}

	weather := models.ReadWeatherDataStruct(data)

	return weather, nil
}

// GetDataOpenWeatherMap call api openweathermap
func GetDataOpenWeatherMap(city string) (models.WeatherData, error) {
	apiConfig, err := common.LoadAPIConfig()
	var data models.WeatherData
	if err != nil {
		return data, err
	}

	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?AppId=" + apiConfig.OpenWeatherMapAPIKey + "&q=" + city)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return data, err
	}

	return data, err
}

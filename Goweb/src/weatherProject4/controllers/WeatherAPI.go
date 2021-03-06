package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"weatherProject4/common"
	"weatherProject4/models"
)

// Add a pool of 5 workers to perform the requests to the external service. If all the workers are busy the originating request should be blocked until one worker is available.
var workersPool chan struct{}
var dataWeather chan string

// WeatherAPI method exposed
func WeatherAPI() {
	workersPool = make(chan struct{}, 5)
	dataWeather = make(chan string)
	http.HandleFunc("/weather/", handlerGetWeatherInfo)
	http.HandleFunc("/scheduler/weather/", handlerSchedulerJob)
	http.ListenAndServe(":8080", nil)
}

func handlerGetWeatherInfo(w http.ResponseWriter, r *http.Request) {
	city := strings.SplitN(r.URL.Path, "/", 3)[2]
	go func() {
		workersPool <- struct{}{}
		defer func() { <-workersPool }()
		data, err := GetWeatherInfo(city)
		if err != nil {
			dataWeather <- err.Error()
			return
		}
		dataWeather <- data
	}()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(<-dataWeather))
}

func handlerSchedulerJob(w http.ResponseWriter, r *http.Request) {
	city := strings.SplitN(r.URL.Path, "/", 4)[3]
	go func() {
		workersPool <- struct{}{}
		defer func() { <-workersPool }()
		go ScheduledJob(city)
	}()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte("Response: 202"))
}

// ScheduledJob check (every 1 hour) of a city and persist it to the database
func ScheduledJob(city string) {
	go func() {
		c := time.Tick(10 * time.Hour)
		for range c {
			GetAndSaveWeatherData(city)
		}
	}()
	select {}
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

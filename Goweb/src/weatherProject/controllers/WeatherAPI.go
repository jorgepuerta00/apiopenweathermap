package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	a "weatherProject/common"

	b "weatherProject/models"
)

// WeatherAPI method exposed
func WeatherAPI() {
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		city := strings.SplitN(r.URL.Path, "/", 3)[2]
		data, err := GetWeatherInfo(city)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(data))
	})
	http.ListenAndServe(":8080", nil)
}

// GetWeatherInfo call api openweathermap
func GetWeatherInfo(city string) (string, error) {
	apiConfig, err := a.LoadAPIConfig()

	if err != nil {
		return err.Error(), err
	}

	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?AppId=" + apiConfig.OpenWeatherMapAPIKey + "&q=" + city)
	if err != nil {
		return err.Error(), err
	}
	defer resp.Body.Close()

	var d b.WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return err.Error(), err
	}

	if d.Message == "" {
		buffer := b.ReadWeatherDataStruct(d)
		return buffer.String(), nil
	}

	message := d.Message
	err = fmt.Errorf("404")

	return message, err
}

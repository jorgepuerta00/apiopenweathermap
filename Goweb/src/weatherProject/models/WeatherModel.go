package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// WeatherItem class
type WeatherItem struct {
	ID          int    `json:"id,omitempty"`
	Main        string `json:"main,omitempty"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
}

// WeatherData class
type WeatherData struct {
	ID         int    `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Visibility int    `json:"visibility,omitempty"`
	Base       string `json:"base,omitempty"`
	Dt         int    `json:"dt,omitempty"`
	Message    string `json:"message,omitempty"`
	Timestamp  string `json:"timestamp,omitempty"`
	Coord      struct {
		Lon float64 `json:"lon,omitempty"`
		Lat float64 `json:"lat,omitempty"`
	} `json:"coord,omitempty"`
	Weather []WeatherItem `json:"weather,omitempty"`
	Main    struct {
		Kelvin   float64 `json:"temp,omitempty"`
		Pressure float64 `json:"pressure,omitempty"`
		Humidity float64 `json:"humidity,omitempty"`
		TempMin  float64 `json:"temp_min,omitempty"`
		TempMax  float64 `json:"temp_max,omitempty"`
	} `json:"main,omitempty"`
	Clouds struct {
		All float64 `json:"all,omitempty"`
	} `json:"clouds,omitempty"`
	Wind struct {
		Speed float64 `json:"speed,omitempty"`
		Deg   float64 `json:"deg,omitempty"`
	} `json:"wind,omitempty"`
	Sys struct {
		ID      int     `json:"id,omitempty"`
		Type    int     `json:"type,omitempty"`
		Message float64 `json:"message,omitempty"`
		Country string  `json:"country,omitempty"`
		Sunrise float64 `json:"sunrise,omitempty"`
		Sunset  float64 `json:"sunset,omitempty"`
	} `json:"sys,omitempty"`
}

// ReadWeatherDataStruct read each item of struct WeatherData
func ReadWeatherDataStruct(data WeatherData) bytes.Buffer {
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(data)
	json.Unmarshal(inrec, &inInterface)

	var buffer bytes.Buffer
	for field, val := range inInterface {
		buffer.WriteString(dumpMap(field, val))
	}

	return buffer
}

func dumpMap(field string, val interface{}) string {
	var detailInformation string
	if !isMessageField(field) {
		detailInformation = field + "\n"
		if isSimpleField(field) {
			detailInformation += "the " + field + " is " + convertValueToString(val) + "\n"
		} else if isWeatherField(field) {
			for i := 0; i < len(val.([]interface{})); i++ {
				nmap := val.([]interface{})[i].(map[string]interface{})
				for key, value := range nmap {
					detailInformation += "the " + key + " is " + convertValueToString(value) + "\n"
				}
			}
		} else {
			nmap := val.(map[string]interface{})
			for key, value := range nmap {
				if isTempField(key) {
					celsiusTemp := celsius(value.(float64))
					detailInformation += "the " + key + " is " + strconv.Itoa(celsiusTemp) + "°C\n"
				} else if isSunField(key) {
					sunTime := time.Unix((int64(value.(float64))), 0)
					detailInformation += "the " + key + " at " + sunTime.Format("Mon, 02 Jan 2006 15:04:05 -0700") + "\n"
				} else {
					detailInformation += "the " + key + " is " + convertValueToString(value) + "\n"
				}
			}
		}
		detailInformation += "\n"
	}
	return detailInformation
}

func fahrenheit(kelvin float64) float64 {
	return (9.0/5.0*(kelvin-273.0) + 32.0)
}

func celsius(kelvin float64) int {
	return int(kelvin - 273.0)
}

func isMessageField(k string) bool {
	var re = regexp.MustCompile(`message|cod`)
	return re.MatchString(k)
}

func isSimpleField(k string) bool {
	var re = regexp.MustCompile(`id|name|visibility|base|dt|timestamp`)
	return re.MatchString(k)
}

func isWeatherField(k string) bool {
	return strings.Contains(k, "weather")
}

func isTempField(k string) bool {
	return strings.Contains(k, "temp")
}

func isSunField(k string) bool {
	return strings.Contains(k, "sun")
}

func convertValueToString(val interface{}) string {
	switch val.(type) {
	case int:
		return strconv.Itoa(int(val.(int)))
	case string:
		return val.(string)
	case float64:
		return fmt.Sprintf("%f", val.(float64))
	default:
		return ""
	}
}

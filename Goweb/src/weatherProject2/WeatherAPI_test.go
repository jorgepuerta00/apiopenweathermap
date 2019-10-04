package main

import (
	"strings"
	"testing"
	a "weatherProject2/controllers"

	. "github.com/smartystreets/goconvey/convey"
)

// TestConveyGetWeatherCali Test function GetWeatherInfo
func TestConveyGetWeather(t *testing.T) {
	t.Parallel()
	city := "cali"
	x, err := a.GetWeatherInfo(city)

	Convey("Starting call api openweathermap", t, func() {
		Convey("When getting response", func() {
			Convey("The data should to cointains "+city+" value in the field Name", func() {
				So(strings.ToUpper(x), ShouldContainSubstring, strings.ToUpper(city))
			})
			Convey("The message should NO be no data found message", func() {
				So(x, ShouldNotEqual, "city not found")
			})
			Convey("The message should be empty", func() {
				So(err, ShouldBeEmpty)
			})
			Convey("The code should NOT be 404", func() {
				So(err, ShouldNotEqual, "404")
			})
		})

	})
}

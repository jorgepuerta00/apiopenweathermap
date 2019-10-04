package common

import (
	"database/sql"
	model "weatherProject2/models"

	// import driver
	_ "github.com/go-sql-driver/mysql"
)

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbName := "globant"
	db, err := sql.Open(dbDriver, dbUser+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

// SelectDataWeather get data from database
func SelectDataWeather(city string) (model.WeatherData, error) {
	var data model.WeatherData
	// do conection
	db := dbConn()
	// Execute the query
	results, err := db.Query("SELECT id, city, country, pressure, humidity, temp, temp_min, temp_max, wind_speed, sunrise, sunset, coord_lon, coord_lat, clouds, weather_main, weather_desc, time_stamp, wind_deg FROM t01weatherinfo WHERE city = ?", city)
	if err != nil {
		return data, err
	}
	for results.Next() {
		data := model.WeatherData{
			Weather: []model.WeatherItem{
				model.WeatherItem{Main: "", Description: ""},
			},
		}
		// for each row, scan the result into our tag composite object
		err = results.Scan(&data.ID, &data.Name, &data.Sys.Country, &data.Main.Pressure, &data.Main.Humidity, &data.Main.Kelvin, &data.Main.TempMin, &data.Main.TempMax, &data.Wind.Speed, &data.Sys.Sunrise, &data.Sys.Sunset, &data.Coord.Lon, &data.Coord.Lat, &data.Clouds.All, &data.Weather[0].Main, &data.Weather[0].Description, &data.Timestamp, &data.Wind.Deg)
		if err != nil {
			return data, err
		}
		return data, err
	}
	defer db.Close()
	return data, nil
}

// InsertDataWeather get data from database
func InsertDataWeather(data model.WeatherData) error {
	// do conection
	db := dbConn()
	// Execute the query
	results, err := db.Query("INSERT INTO t01weatherinfo (id, city, country, pressure, humidity, temp, temp_min, temp_max, wind_speed, sunrise, sunset, coord_lon, coord_lat, clouds, weather_main, weather_desc, wind_deg) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", &data.ID, &data.Name, &data.Sys.Country, &data.Main.Pressure, &data.Main.Humidity, &data.Main.Kelvin, &data.Main.TempMin, &data.Main.TempMax, &data.Wind.Speed, &data.Sys.Sunrise, &data.Sys.Sunset, &data.Coord.Lon, &data.Coord.Lat, &data.Clouds.All, &data.Weather[0].Main, &data.Weather[0].Description, &data.Wind.Deg)
	if err != nil {
		return err
	}
	defer results.Close()
	defer db.Close()
	return nil
}

// DeleteDataWeather get data from database
func DeleteDataWeather(data model.WeatherData) error {
	// do conection
	db := dbConn()
	// Execute the query
	results, err := db.Query("DELETE FROM t01weatherinfo WHERE city = ?", &data.Name)
	if err != nil {
		return err
	}
	defer results.Close()
	defer db.Close()
	return nil
}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const weatherAPIKey = "5bc613c75578a94e31e2f63f5757caa5"
const apiKey2 = "3abc5ce17df4a5bca8476c92ae805921"

type Coordinates struct {
	City string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

type Weather struct {
	Timezone string         `json:"timezone"`
	Current  CurrentWeather `json:"current"`
	Alerts   []WeatherAlert `json:"alerts"`
}

type CurrentWeather struct {
	Dt         int     `json:"dt"`
	Sunrise    int     `json:"sunrise"`
	Sunset     int     `json:"sunset"`
	Temp       float64 `json:"temp"`
	FeelsLike  float64 `json:"feels_like"`
	Pressure   int     `json:"pressure"`
	Humidity   int     `json:"humidity"`
	DewPoint   float64 `json:"dew_point"`
	Uvi        float64 `json:"uvi"`
	Clouds     int     `json:"clouds"`
	Visibility int     `json:"visibility"`
	WindSpeed  float64 `json:"wind_speed"`
	WindDeg    int     `json:"wind_deg"`
	WindGust   float64 `json:"wind_gust"`
	Extra      []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
}

type WeatherAlert struct {
	SenderName  string   `json:"sender_name"`
	Event       string   `json:"event"`
	Start       int64    `json:"start"`
	End         int64    `json:"end"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

func getCoordinatesByZip(zipCode string) (Coordinates, error) {

	apiURL := fmt.Sprintf("https://api.openweathermap.org/geo/1.0/zip?zip=%s,US&appid=%s", zipCode, weatherAPIKey)

	res, err := http.Get(apiURL)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Coordinates{}, err
	}

	var locResp Coordinates
	if err := json.Unmarshal(body, &locResp); err != nil {
		return Coordinates{}, err
	}

	return locResp, nil
}

func getWeather(longitude float64, latitude float64) (Weather, error) {
	apiURL := fmt.Sprintf("https://api.openweathermap.org/data/3.0/onecall?lat=%f&lon=%f&exclude=minutely,hourly,daily,lat,lon&units=imperial&appid=%s", latitude, longitude, weatherAPIKey)

	res, err := http.Get(apiURL)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("body: %s\n", res.Status)
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Weather{}, err
	}
	var weatherResp Weather
	if err := json.Unmarshal(body, &weatherResp); err != nil {
		return Weather{}, err
	}

	return weatherResp, nil
}

func main() {
	zip := flag.String("zip", "", "Specify a zip code")

	flag.Parse()

	if *zip == "" {
		fmt.Println("Please provide a zip code using the -zip flag.")
		os.Exit(1)
	}

	locData, err := getCoordinatesByZip(*zip)
	if err != nil {
		log.Fatal("Error:", err)
	}
	fmt.Printf("City: %s| latitude: %f, longitude: %f", locData.City, locData.Lat, locData.Lon)

	weatherData, err := getWeather(locData.Lon, locData.Lat)

	if err != nil {
		log.Fatal("Error:", err)
	}
	fmt.Printf("Current temperature: %f\n", weatherData.Current.Temp)

}

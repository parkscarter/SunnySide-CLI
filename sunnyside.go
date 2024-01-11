package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"unicode"
)

const weatherAPIKey = "5bc613c75578a94e31e2f63f5757caa5"

/*
This struct represents the response given when searching a location based on zip code or name
*/
type Coordinates struct {
	City string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

// This simply represents an array of Coordinates objects
type CoordinatesArray []Coordinates

/*
This struct represents the response given when searching the weather based on latitude and logitude
This struct also contains a struct (CurrentWeather), and an array structs (WeatherAlert) which are nested in the response
*/
type Weather struct {
	Timezone string         `json:"timezone"`
	Current  CurrentWeather `json:"current"`
	Alerts   []WeatherAlert `json:"alerts"`
}

/*
This struct represents information abouut the current weather, which is nested in the weather struct
*/
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

/*
This struct is also nested within the larger Weather struct, and represents information regarding warnings
*/
type WeatherAlert struct {
	SenderName  string   `json:"sender_name"`
	Event       string   `json:"event"`
	Start       int64    `json:"start"`
	End         int64    `json:"end"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

func getWindDir(degrees float64) string {
	cardinalDirections := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW", "N"}

	index := int((degrees + 22.5) / 45)
	if index < 0 {
		index = 0
	} else if index > 7 {
		index = 0
	}

	return cardinalDirections[index]
}

func getCloudCoverage(percentage int) string {
	switch {
	case percentage <= 10:
		return "Sunny"
	case percentage <= 50:
		return "Partly Cloudy"
	case percentage <= 80:
		return "Mostly Cloudy"
	default:
		return "Very Cloudy"
	}
}

/*
This function takes longitude and latitude as parameters, and returns an instance of the Weather struct
*/
func getWeather(longitude float64, latitude float64) {
	apiURL := fmt.Sprintf("https://api.openweathermap.org/data/3.0/onecall?lat=%f&lon=%f&exclude=minutely,hourly,daily,lat,lon&units=imperial&appid=%s", latitude, longitude, weatherAPIKey)

	res, err := http.Get(apiURL)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	var weatherResp Weather
	if err := json.Unmarshal(body, &weatherResp); err != nil {
		return
	}
	tempInt := int(math.Round(weatherResp.Current.Temp))
	feelsInt := int(math.Round(weatherResp.Current.FeelsLike))
	windInt := int(math.Round(weatherResp.Current.WindSpeed))
	windDir := getWindDir(float64(weatherResp.Current.WindDeg))
	cloudCov := getCloudCoverage(weatherResp.Current.Clouds)

	numConditions := len(weatherResp.Current.Extra)

	conditions := ""

	for i := 0; i < numConditions; i++ {
		conditions += weatherResp.Current.Extra[i].Description
		if i+1 < len(weatherResp.Current.Extra) {
			conditions += ", "
		}
	}

	numAlerts := len(weatherResp.Alerts)

	alerts := ""

	for i := 0; i < numAlerts; i++ {
		alerts += weatherResp.Alerts[i].Event
		if i+1 < len(weatherResp.Alerts) {
			alerts += "\n- "
		}
	}

	fmt.Printf("Current temperature: %d    Feels like: %d\n\n", tempInt, feelsInt)

	fmt.Printf("Wind Conditions: %d MPH from the %s\n\n", windInt, windDir)

	fmt.Printf("%s with %s\n\n", cloudCov, conditions)

	fmt.Printf("Alerts:\n- %s", alerts)

	return
}

/*
This function is called to return coordinates based on zip code (US)
Utelizes external api (OpenWeatherMap)
*/
func getCoordinatesByZip(zipCode string) {

	apiURL := fmt.Sprintf("https://api.openweathermap.org/geo/1.0/zip?zip=%s,US&appid=%s", zipCode, weatherAPIKey)

	res, err := http.Get(apiURL)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	var locResp Coordinates
	if err := json.Unmarshal(body, &locResp); err != nil {
		return
	}
	fmt.Printf("\n\n\n\nLocal weather in %s:\n\n", locResp.City)
	getWeather(locResp.Lon, locResp.Lat)

	return
}

func takeZipInput() {
	var input string
	fmt.Print("\n\n\n\nEnter 'q' to quit or 'b' to go back\n\n")
	fmt.Printf("Enter a US based Zip code:  \n")
	for {
		_, err := fmt.Scan(&input)

		if input == "q" {
			os.Exit(1)
		}

		if input == "b" {
			return
		}

		//Check if zip code is numerical
		isNumeric := true
		for _, char := range input {
			if !unicode.IsDigit(char) {
				isNumeric = false
				break
			}
		}

		//Check if the zip code's length is 5
		if err == nil && len(input) == 5 && isNumeric == true {
			getCoordinatesByZip(input)
			break
		}
		fmt.Println("Please enter a valid zip code:")
	}
	return
}

func getCoordinatesByCity(city string, state string, country string) {
	apiURL := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s,%s,%s&limit=1&appid=%s", city, state, country, weatherAPIKey)

	res, err := http.Get(apiURL)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Printf("error converting to body: %s\n", err)
		return
	}

	var locRespArray CoordinatesArray
	if err := json.Unmarshal(body, &locRespArray); err != nil {
		fmt.Printf("error unmarshalling: %s\n", err)
		return
	}

	if len(locRespArray) == 0 {
		fmt.Println("No coordinates found.")
		return
	}

	locResp := locRespArray[0]

	fmt.Printf("\n\nLocal weather in %s:\n\n", locResp.City)
	getWeather(locResp.Lon, locResp.Lat)

	return
}

func takeCityInput() {
	var countryInput string
	var stateInput string
	var cityInput string

	fmt.Print("\n\n\n\nEnter 'q' to quit or 'b' to go back\n\n")
	fmt.Printf("Enter a country code (ex. US):\n")

	for {
		_, err := fmt.Scan(&countryInput)

		if countryInput == "q" {
			os.Exit(1)
		}

		if countryInput == "b" {
			return
		}

		isChars := true
		for _, char := range countryInput {
			if unicode.IsDigit(char) {
				isChars = false
				break
			}
		}

		if err == nil && len(countryInput) == 2 && isChars == true {
			break
		}
		fmt.Println("Please enter a valid country code (ex. US) (enter 'q' to quit):")
	}

	fmt.Print("\n\n\n\nEnter 'q' to quit or 'b' to go back\n\n")
	fmt.Printf("Enter a state code (ex: IA):\n")

	for {
		_, err := fmt.Scan(&stateInput)

		if stateInput == "q" {
			os.Exit(1)
		}

		if stateInput == "b" {
			return
		}

		isChars := true
		for _, char := range stateInput {
			if unicode.IsDigit(char) {
				isChars = false
				break
			}
		}

		if err == nil && len(stateInput) == 2 && isChars == true {
			break
		}
		fmt.Println("Please enter a valid state code (ex. IA) (enter 'q' to quit):")
	}

	fmt.Print("\n\n\n\nEnter 'q' to quit or 'b' to go back\n\n")
	fmt.Printf("Enter the name of a city (ex. Oskaloosa):\n")

	for {
		_, err := fmt.Scan(&cityInput)

		if cityInput == "q" {
			os.Exit(1)
		}

		if cityInput == "b" {
			return
		}

		isChars := true
		for _, char := range cityInput {
			if unicode.IsDigit(char) {
				isChars = false
				break
			}
		}

		if err == nil && isChars == true {
			break
		}
		fmt.Println("Please enter a valid City (ex. Oskaloosa) (enter 'q' to quit):")
	}
	getCoordinatesByCity(cityInput, stateInput, countryInput)

}

func main() {
	var input string
	fmt.Printf("\n\n\n\n_______________________________________________________________________________\n\n")
	fmt.Printf("Hello and welcome to SunnySide Weather! This interactive CLI was built as practice;\nhowever it still has some neat functionality to explore; most importantly, \nthis program will return the current weather at a location specified by the user\n\n")
	fmt.Printf("Enter 'z' to search by zipcode, 'l' to search by city name, or 'q' to quit\n")
	for {
		_, err := fmt.Scan(&input)
		if err != nil {
			fmt.Println("Error reading input:", err)
		} else if input == "q" {
			return
		} else if input == "z" {
			takeZipInput()
		} else if input == "l" {
			takeCityInput()
		}
		fmt.Printf("\n\nEnter 'z' to search by zipcode, 'l' to search by city name, or 'q' to quit\n")
	}

}

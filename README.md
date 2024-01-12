# SunnySide-CLI
This CLI was written as my first real practice with Go. It makes calls to an external API (OpenWeatherMap) to gather information about the current weather based on location. 

Locations can be entered via zip code or by city name (which requires country and state codes)

Although the API returns an abundance of information, I decided to filter out most of it so more focus is drawn to important information like temperature, real feel, wind speed, wind direction, and any potential warnings for the area. A future implementation could include an option to display more information, as well as a forecast.

To run the program:

First download the project and install Go on your machine if it isn't installed already

Next, use your command line to cd to where the file is located

Lastly, enter command "go run sunnyside.go" and follow the directions in your terminal

package app

import (
	"WeatherApi/weather"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func defaultRequest(w http.ResponseWriter, request *http.Request) {
	mUrl := request.RequestURI

	if strings.Contains(mUrl, "icon") || strings.Contains(mUrl, "fav") {
		return
	}

	fmt.Printf("Also called here, with URL \"%s\"\n", request.RequestURI)

	w.WriteHeader(http.StatusPartialContent)
	w.Header().Set("Content-Type", "application/json")

	mkErr := mkError(206, "Supported API Connections: \\\"weather\"\\")
	jsonData, err := json.Marshal(mkErr)
	if err != nil {
		fmt.Printf("JSON Marshal Error: %v\n", err)
		return
	}

	defer w.Write(jsonData)
}

func weatherRequest(writer http.ResponseWriter, request *http.Request) {
	reqUrl := strings.TrimSpace(request.RequestURI)
	fmt.Printf("Received URL (\"%s\")\n", reqUrl)

	if reqUrl == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	mUrl, err := url.ParseRequestURI(reqUrl)
	if err != nil {
		fmt.Printf("ParseRequestURI_Error: %v\n", err)
	}

	sUrl := mUrl.Path[len("/weather/"):]
	fmt.Printf("S-Url: \"%s\"\n", sUrl)

	params := mUrl.Query()
	var mp map[string]string
	if params != nil && len(params) != 0 {
		mp, _ = retrieveParams(params)
	}

	timezone := ""
	var lat, lon *float64

	if mp != nil && len(mp) != 0 {
		timezone = mp["tz"]

		if len(mp["lat"]) != 0 {
			fl, err := convertStringToFloat(mp["lat"])
			if err != nil {
				fmt.Printf("LocationParserError(Latitude) :: %v\n", err)
			}

			lat = &fl
		}

		if len(mp["lon"]) != 0 {
			fl, err := convertStringToFloat(mp["lat"])
			if err != nil {
				fmt.Printf("LocationParserError(Longitude) :: %v\n", err)
			}

			lon = &fl
		}
	}

	if len(sUrl) == 0 || sUrl == "" {
		defer writer.Write([]byte("{\"error\": 206, \"message\": \"Must contain one of these give URL paths: current, today, tomorrow, week, day\"}"))
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusPartialContent)

		return
	}

	if lat == nil || lon == nil {
		defer writer.Write([]byte("{\"error\": 400, \"message\": \"Latitude (lat) and/or Longitude (lon) URL parameters are missing\"}"))
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	// Check for common requests
	fmt.Printf("URL Substring: %s\n", sUrl)
	switch sUrl {
	case "current":
		// Retrieve the current weather for the current timestamp
		w, err := weather.Current(timezone, *lat, *lon)

		if err != nil {
			fmt.Printf("CurrentWeatherError :: %v\n", err)
		}

		fmt.Printf("Current Weather: %v\n", w)

		b, e := json.Marshal(w)
		if e != nil {
			fmt.Printf("JsonError(CurrentWeather) :: %v\n", e)
		}

		defer writer.Write(b)
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
	case "today":
		// Check for the weather for the current day
		w, err := weather.SpecificDay(time.Now(), timezone, *lat, *lon)

		if err != nil {
			fmt.Printf("WeatherTodayError :: %v\n", err)
		}

		b, e := json.Marshal(w.Weather)
		if e != nil {
			fmt.Printf("JsonError(WeatherToday) :: %v\n", e)
		}

		defer writer.Write(b)
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
	case "tomorrow":
		// Check for the weather for tomorrow
		w, err := weather.SpecificDay(time.Now().Add(time.Hour*24), timezone, *lat, *lon)

		if err != nil {
			fmt.Printf("WeatherTodayError :: %v\n", err)
		}

		b, e := json.Marshal(w.Weather)
		if e != nil {
			fmt.Printf("JsonError(WeatherToday) :: %v\n", e)
		}

		defer writer.Write(b)
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
	case "week":
		// Check for the weather for this week, looping through all days on the current week (Monday - Sunday)
	case "day":
		// Check for the weather for a specific day
	default:
		// Unknown URL, throw 400 Bad Request
		writer.Header().Set("Content-Type", "application/json")
		defer writer.Write([]byte(fmt.Sprintf("{\"error\":404,\"message\":\"Unknown URL, received \\\"%s\\\"\"}", request.RequestURI)))
		if err != nil {
			return
		}
		writer.WriteHeader(http.StatusNotFound)
	}
}

func RunServer(port int) error {
	sHandle := http.NewServeMux()

	sHandle.HandleFunc("/weather/", weatherRequest)
	sHandle.HandleFunc("/", defaultRequest)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), sHandle)
}

package app

import (
	"WeatherApi/weather"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func defaultRequest(writer http.ResponseWriter, request *http.Request) {
	fmt.Printf("Also called here, with URL \"%s\"\n", request.RequestURI)
}

var multiplierCount time.Duration

func weatherRequest(writer http.ResponseWriter, request *http.Request) {
	// reqUrl := strings.TrimSpace(request.RequestURI[len("/weather/"):])
	reqUrl := strings.TrimSpace(request.RequestURI)
	fmt.Printf("Received URL \"%s\"\n", reqUrl)

	if reqUrl == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	mUrl, err := url.ParseRequestURI(reqUrl)
	if err != nil {
		fmt.Printf("ParseRequestURI_Error: %v\n", err)
	}

	params := mUrl.Query()

	timezone := func() string {
		if params.Has("tz") {
			return params.Get("tz")
		}

		if params.Has("timezone") {
			return params.Get("timezone")
		}

		return ""
	}()

	lat := func() float64 {
		if params.Has("lat") {
			fl, _ := strconv.ParseFloat(params.Get("lat"), 64)
			return fl
		}

		if params.Has("latitude") {
			fl, _ := strconv.ParseFloat(params.Get("latitude"), 64)
			return fl
		}

		return 0.0
	}()

	lon := func() float64 {
		if params.Has("lon") {
			fl, _ := strconv.ParseFloat(params.Get("lon"), 64)
			return fl
		}

		if params.Has("longitude") {
			fl, _ := strconv.ParseFloat(params.Get("longitude"), 64)
			return fl
		}

		return 0.0
	}()

	// Check for common requests
	switch mUrl.Path[len("/weather/"):] {
	case "current":
		// Retrieve the current weather for the current timestamp
	case "today":
		// Check for the weather for the current day
		w, err := weather.SpecificDay(time.Now(), timezone, lat, lon)

		if err != nil {
			fmt.Printf("[ERROR:\"/weather/today\"] %v\n", err)
		}

		fmt.Printf("Weather: %v\n", w)
	case "tomorrow":
		// Check for the weather for tomorrow
	case "week":
		// Check for the weather for this week, looping through all days on the current week (Monday - Sunday)
	default:
		// Unknown URL, throw 400 Bad Request
		if strings.HasPrefix(reqUrl, "day") {
			break
		}

		writer.Header().Set("Content-Type", "application/json")
		defer writer.Write([]byte("{\"error\":400,\"message\":\"Unknown URL\"}"))
		if err != nil {
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func RunServer(port int) error {
	sHandle := http.NewServeMux()

	sHandle.HandleFunc("/weather/", weatherRequest)
	sHandle.HandleFunc("/", defaultRequest)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), sHandle)
	if err != nil {
		return err
	}

	return nil
}

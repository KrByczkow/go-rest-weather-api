package app

import (
	"WeatherApi/weather"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

func initializeClientRequest(w http.ResponseWriter, r *http.Request) (*time.Time, string, *float64, *float64, error) {
	params := mkMapParams(r)
	if params == nil || len(params) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		bytes, err := json.Marshal(mkError(400, "Requires either the city and country parameters, or the latitude along with longitude"))
		if err != nil {
			fmt.Printf("JsonBuilderError: %v\n", err)
		}

		defer w.Write(bytes)
	}

	var tz string
	var lat, lon *float64
	var tm *time.Time

	if params != nil && len(params) != 0 {
		tz = params["tz"]

		if len(params["lat"]) != 0 {
			fl, err := convertStringToFloat(params["lat"])
			if err != nil {
				fmt.Printf("LocationParserError(Latitude) :: %v\n", err)
			}

			lat = &fl
		}

		if len(params["lon"]) != 0 {
			fl, err := convertStringToFloat(params["lon"])
			if err != nil {
				fmt.Printf("LocationParserError(Longitude) :: %v\n", err)
			}

			lon = &fl
		}

		if len(params["tm"]) != 0 {
			retTime, err := time.Parse("2006-01-02", params["tm"])
			if err != nil {
				fmt.Printf("TimeParseError :: %v\n", err)
			}

			tm = &retTime
		}
	}

	w.Header().Set("Content-Type", "application/json")

	return tm, tz, lat, lon, nil
}

func currentWeatherRequest(w http.ResponseWriter, r *http.Request) {
	_, tz, lat, lon, err := initializeClientRequest(w, r)
	if err != nil {
		fmt.Printf("{RequestError} %v\n", err)
		return
	}

	wt, err := weather.Current(tz, *lat, *lon)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		errBody := mkError(500, fmt.Sprintf("Cannot obtain weather due to this following error: %v", err))
		errJson, err := json.Marshal(errBody)
		if err != nil {
			fmt.Printf("Error making JSON Body: %v\n", err)
			return
		}

		w.Write(errJson)
		return
	}

	jsonBody, err := json.Marshal(wt)
	if err != nil {
		fmt.Printf("Error making JSON Body: %v\n", err)
		return
	}

	w.WriteHeader(200)
	w.Write(jsonBody)
}

func todayWeatherRequest(w http.ResponseWriter, r *http.Request) {
	_, tz, lat, lon, err := initializeClientRequest(w, r)
	if err != nil {
		fmt.Printf("{RequestError} %v\n", err)
		return
	}

	wt, err := weather.SpecificDay(time.Now(), tz, *lat, *lon)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		errBody := mkError(500, fmt.Sprintf("Cannot obtain weather due to this following error: %v", err))
		errJson, err := json.Marshal(errBody)
		if err != nil {
			fmt.Printf("Error making JSON Body: %v\n", err)
			return
		}

		w.Write(errJson)
		return
	}

	jsonBody, err := json.Marshal(wt)
	if err != nil {
		fmt.Printf("Error making JSON Body: %v\n", err)
		return
	}

	w.WriteHeader(200)
	w.Write(jsonBody)
}

func dayWeatherRequest(w http.ResponseWriter, r *http.Request) {
	tm, tz, lat, lon, err := initializeClientRequest(w, r)
	if err != nil {
		fmt.Printf("{RequestError} %v\n", err)
		return
	}

	if tm == nil {
		w.WriteHeader(http.StatusBadRequest)

		errBody := mkError(http.StatusBadRequest, fmt.Sprintf("The 'date' parameter is missing"))
		errJson, err := json.Marshal(errBody)

		if err != nil {
			fmt.Printf("Error making JSON Body: %v\n", err)
			return
		}

		w.Write(errJson)
		return
	}

	wt, err := weather.SpecificDay(*tm, tz, *lat, *lon)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		errBody := mkError(500, fmt.Sprintf("Cannot obtain weather due to this following error: %v", err))
		errJson, err := json.Marshal(errBody)
		if err != nil {
			fmt.Printf("Error making JSON Body: %v\n", err)
			return
		}

		w.Write(errJson)
		return
	}

	jsonBody, err := json.Marshal(wt)
	if err != nil {
		fmt.Printf("Error making JSON Body: %v\n", err)
		return
	}

	w.WriteHeader(200)
	w.Write(jsonBody)
}

func ExtRunServer(port int) error {
	rt := chi.NewRouter()

	rt.Get("/weather/current", currentWeatherRequest)
	rt.Get("/weather/today", todayWeatherRequest)
	rt.Get("/weather/day", dayWeatherRequest)

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: rt,

		ReadTimeout:       time.Second * 30,
		WriteTimeout:      time.Second * 30,
		IdleTimeout:       time.Second * 60,
		ReadHeaderTimeout: time.Second * 15,
	}

	return srv.ListenAndServe()
}

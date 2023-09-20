package app

import (
	"WeatherApi/weather"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

func initializeClientRequest(w http.ResponseWriter, r *http.Request) (*time.Time, string, *float64, *float64, error) {
	w.Header().Set("Content-Type", "application/json")

	params, paramCount := mkMapParams(r)
	if params == nil || len(params) == 0 {
		if paramCount == 0 {
			return nil, "", nil, nil, errors.New("no parameters have been given to the URL request")
		}

		return nil, "", nil, nil, errors.New("invalid parameters have been given to the URL request")
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

	if lat == nil || lon == nil {
		return tm, tz, nil, nil, errors.New("latitude and/or longitude parameters are missing")
	}

	return tm, tz, lat, lon, nil
}

func currentWeatherRequest(w http.ResponseWriter, r *http.Request) {
	_, tz, lat, lon, err := initializeClientRequest(w, r)
	if err != nil {
		errTag := mkError(http.StatusBadRequest, err.Error())

		if iErr := sendData(w, errTag); iErr != nil {
			fmt.Printf("[ResponseSendError] %v\n", iErr)
		}

		return
	}

	wt, err := weather.Current(tz, *lat, *lon)
	if err != nil {
		errBody := mkError(http.StatusInternalServerError, fmt.Sprintf("Cannot obtain weather due to this following error: %v", err))

		if iErr := sendData(w, errBody); iErr != nil {
			fmt.Printf("[ResponseSendError] %v\n", err)
		}

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
		errTag := mkError(http.StatusBadRequest, fmt.Sprintf("request error (%v)", err))

		if nErr := sendData(w, errTag); nErr != nil {
			fmt.Printf("[ResponseSendError] %v\n", nErr)
		}

		return
	}

	wt, err := weather.SpecificDay(time.Now(), tz, *lat, *lon)
	if err != nil {
		errBody := mkError(500, fmt.Sprintf("Cannot obtain weather due to this following error: %v", err))

		if iErr := sendData(w, errBody); iErr != nil {
			fmt.Printf("[ResponseSendError:InternalError] %v\n", err)
		}

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
		errTag := mkError(http.StatusBadRequest, fmt.Sprintf("request error (%v)", err))

		if nErr := sendData(w, errTag); nErr != nil {
			fmt.Printf("[ResponseSendError] %v\n", nErr)
		}

		return
	}

	if tm == nil {
		errBody := mkError(http.StatusBadRequest, fmt.Sprintf("The 'date' parameter is missing"))

		if nErr := sendData(w, errBody); nErr != nil {
			fmt.Printf("[ResponseSendError] %v\n", nErr)
		}

		return
	}

	wt, err := weather.SpecificDay(*tm, tz, *lat, *lon)
	if err != nil {
		errBody := mkError(http.StatusInternalServerError, fmt.Sprintf("Cannot obtain weather due to this following error: %v", err))

		if iErr := sendData(w, errBody); iErr != nil {
			fmt.Printf("[ResponseSendError] %v\n", iErr)
		}

		return
	}

	jsonBody, err := json.Marshal(wt)
	if err != nil {
		fmt.Printf("Error making JSON Body: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBody)
}

func ExtRunServer(bindAddr string) error {
	rt := chi.NewRouter()

	rt.Get("/weather/current", currentWeatherRequest)
	rt.Get("/weather/today", todayWeatherRequest)
	rt.Get("/weather/day", dayWeatherRequest)

	srv := http.Server{
		Addr:    bindAddr,
		Handler: rt,

		ReadTimeout:       time.Second * 30,
		WriteTimeout:      time.Second * 30,
		IdleTimeout:       time.Second * 60,
		ReadHeaderTimeout: time.Second * 15,
	}

	return srv.ListenAndServe()
}

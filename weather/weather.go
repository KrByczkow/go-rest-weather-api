package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

func openRequest(strUrl, tz string, tm time.Time, lat, lon float64) (string, error) {
	params := url.Values{}
	params.Set("lat", fmt.Sprintf("%f", lat))
	params.Set("lon", fmt.Sprintf("%f", lon))

	if !tm.IsZero() {
		params.Set("date", fmt.Sprintf("%04d-%02d-%02d", tm.Year(), tm.Month(), tm.Day()))
	}

	if tz != "" {
		params.Set("tz", tz)
	}

	reqUrl := fmt.Sprintf("%s?%s", strUrl, params.Encode())
	fmt.Printf("Making a GET Request to %s\n", reqUrl)
	response, err := http.Get(reqUrl)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	return string(body), nil
}

func jsonCurrent(tz string, lat, lon float64) (string, error) {
	return openRequest(CurrentWeatherUrl, tz, time.Time{}, lat, lon)
}

func jsonDay(tm time.Time, tz string, lat, lon float64) (string, error) {
	return openRequest(HourlyWeatherUrl, tz, tm, lat, lon)
}

func Current(tz string, lat, lon float64) (CurrentWeather, error) {
	jsonData, err := jsonCurrent(tz, lat, lon)
	if err != nil {
		return CurrentWeather{}, err
	}

	var weather CurrentWeather
	err = json.Unmarshal([]byte(jsonData), &weather)

	if err != nil {
		return CurrentWeather{}, err
	}

	return weather, nil
}

func SpecificDay(tm time.Time, tz string, lat, lon float64) (DayWeather, error) {
	jsonData, err := jsonDay(tm, tz, lat, lon)
	if err != nil {
		return DayWeather{}, err
	}

	var weather DayWeather
	err = json.Unmarshal([]byte(jsonData), &weather)

	if err != nil {
		return DayWeather{}, err
	}

	return weather, nil
}

package app

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func retrieveParams(params url.Values) map[string]string {
	var outMap = map[string]string{}

	for k, v := range params {
		value := strings.Join(v, ", ")

		if strings.HasPrefix(k, "lat") {
			if _, ok := outMap["lat"]; !ok {
				outMap["lat"] = value
			}
		}

		if strings.HasPrefix(k, "lon") {
			if _, ok := outMap["lon"]; !ok {
				outMap["lon"] = value
			}
		}

		if k == "tz" || k == "timezone" {
			if _, ok := outMap["tz"]; !ok {
				outMap["tz"] = value
			}
		}

		if k == "date" || k == "time" {
			if _, ok := outMap["tm"]; !ok {
				outMap["tm"] = value
			}
		}
	}

	return outMap
}

func mkMapParams(r *http.Request) map[string]string {
	mUrl, err := url.ParseRequestURI(r.RequestURI)
	if err != nil {
		fmt.Printf("ParseRequestUriError: %v\n", err)
		return nil
	}

	params := mUrl.Query()
	var mp map[string]string

	if params != nil && len(params) != 0 {
		mp = retrieveParams(params)
	}

	return mp
}

func mkError(errCode int, message string) ErrorMessage {
	return ErrorMessage{errCode, message}
}

func convertStringToFloat(str string) (float64, error) {
	if str == "" || strings.TrimSpace(str) == "" {
		return 0, fmt.Errorf("string is empty")
	}

	fl, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, err
	}

	return fl, nil
}

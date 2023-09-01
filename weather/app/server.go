package app

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func defaultRequest(writer http.ResponseWriter, request *http.Request) {
	fmt.Printf("Also called here, with URL \"%s\"\n", request.RequestURI)
}

func weatherRequest(writer http.ResponseWriter, request *http.Request) {
	reqUrl := strings.TrimSpace(request.RequestURI[len("/weather/"):])
	fmt.Printf("Received URL \"%s\"\n", reqUrl)

	if reqUrl == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	escUrl, err := url.QueryUnescape(reqUrl)
	if err != nil {
		fmt.Printf("Error unescaping: %v\n", err)
	}

	reqUrl = strings.TrimSpace(escUrl)
	fmt.Printf("New URL: \"%s\"\n", reqUrl)

	// Check for common requests
	switch reqUrl {
	case "current":
		break
	case "today":
		// Check for the weather for the current day
		break
	case "tomorrow":
		// Check for the weather for tomorrow
		break
	case "week":
		// Check for the weather for this week, looping through all days on the current week (Monday - Sunday)
		break
	default:
		// Unknown URL, throw 400 Bad Request
		if strings.HasPrefix(reqUrl, "day") {
			break
		}

		writer.WriteHeader(http.StatusBadRequest)

		break
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

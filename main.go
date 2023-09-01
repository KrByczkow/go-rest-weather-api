package main

import (
	"WeatherApi/weather"
	"WeatherApi/weather/app"
	"fmt"
	"time"
)

func main() {
	current, err := weather.Current("Europe/Berlin", 52.52, 13.4)
	if err != nil {
		fmt.Printf("[Error:Main] %v\n", err)
		return
	}

	day, err := weather.SpecificDay(time.Now(), "Europe/Berlin", 52.52, 13.4)
	if err != nil {
		fmt.Printf("[Error:Main] %v\n", err)
		return
	}

	tomorrow, err := weather.SpecificDay(time.Now().Add(time.Hour*24), "Europe/Berlin", 52.52, 13.4)
	if err != nil {
		fmt.Printf("[Error:Main] %v\n", err)
		return
	}

	yesterday, err := weather.SpecificDay(time.Now(), "Europe/Berlin", 52.52, 13.4)
	if err != nil {
		fmt.Printf("[Error:Main] %v\n", err)
		return
	}

	fmt.Printf("Current: %v\n\n", current.Weather)
	fmt.Printf("Today: %v\n\n", day.Weather)
	fmt.Printf("Tomorrow: %v\n\n", tomorrow.Weather)
	fmt.Printf("Yesterday: %v\n\n", yesterday.Weather)

	err = app.RunServer(8080)
	if err != nil {
		fmt.Printf("[Error:Main] %v\n", err)
	}
}

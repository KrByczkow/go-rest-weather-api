package main

import (
	"WeatherApi/weather/app"
	"fmt"
)

func main() {
	err := app.ExtRunServer(8080)
	if err != nil {
		fmt.Printf("[Error:Main] %v\n", err)
	}
}

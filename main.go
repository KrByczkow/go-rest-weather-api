package main

import (
	"WeatherApi/weather/app"
	"flag"
	"fmt"
)

func main() {
	var bindAddr string

	flag.StringVar(&bindAddr, "b", ":8080", "Bind Address to the Webserver")
	flag.Parse()

	fmt.Printf("Running webserver on bind address %s\n", bindAddr)
	err := app.ExtRunServer(bindAddr)
	if err != nil {
		fmt.Printf("[Error:Main] %v\n", err)
	}
}

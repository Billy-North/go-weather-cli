package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gosuri/uilive"
)

func getUserTextInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func getUserInputBoolean() bool {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	switch scanner.Text() {
	case "y":
		return true
	case "n":
		return false
	}

	fmt.Println("Invalid input - select one of y/n")
	return getUserInputBoolean()
}

func getUserNumberInput() int {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	s := scanner.Text()

	parsedInt, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		fmt.Println("Input a valid number")
		return getUserNumberInput()
	}

	//Should be safe as bit size 0 has been selected from ParseInt
	return int(parsedInt)

}

func displaySelectLocation(results []OpenMetoSearchNameResult) OpenMetoSearchNameResult {
	if len(results) == 0 {
		fmt.Println("No results matching this search were found please try again.")
		os.Exit(0)
	}

	fmt.Println("Confirm the location to check weather")
	for idx, e := range results {
		fmt.Printf("%d. %s - %s  \n", idx, e.Name, e.Country)
	}

	selectedValue := getUserNumberInput()
	if selectedValue >= len(results) || selectedValue < 0 {
		fmt.Println("Invalid value - please select a from a value displayed")
		return displaySelectLocation(results)
	}

	return results[selectedValue]
}

func formatWeatherData(locationName string, weatherData OpenMetoCurrentWeatherResult) string {
	return fmt.Sprintf(
		"The Current Weather in %v \nTemperature: %v°C\nWindspeed: %vkm/h\nWinddirection: %v°\n",
		locationName, weatherData.Temperature, weatherData.Windspeed, weatherData.Winddirection)
}

func userSelectUpdate() bool {
	fmt.Print("Would you like to recieve live updates? - enter y/n: ")
	return getUserInputBoolean()
}

func stopWriter(writer *uilive.Writer) {
	fmt.Println("Stopping Writer...")
	writer.Stop()
}

func pollWeatherDataLoop(selectedLocation OpenMetoSearchNameResult) {
	const POLL_TIMEOUT_SECONDS = 60
	writer := uilive.New()
	writer.Start()
	for true {
		weatherData, err := QueryCurrentWeather(
			LocationCoordinates{
				Latitude:  selectedLocation.Latitude,
				Longitude: selectedLocation.Longitude,
			})

		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}

		fmt.Fprintf(writer, formatWeatherData(selectedLocation.Name, weatherData))
		time.Sleep(time.Second * POLL_TIMEOUT_SECONDS)
	}
}

func displayStaticWeather(selectedLocation OpenMetoSearchNameResult) {
	weatherData, err := QueryCurrentWeather(
		LocationCoordinates{
			Latitude:  selectedLocation.Latitude,
			Longitude: selectedLocation.Longitude,
		})

	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(formatWeatherData(selectedLocation.Name, weatherData))
}

func main() {
	fmt.Println("Welcome to the Weather Checker...")
	fmt.Print("Enter a location: ")
	inputLocation := getUserTextInput()
	fmt.Println("The selected location is", inputLocation)
	avaliableLocations, err := QuerySearchLocations(inputLocation)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	selectedLocation := displaySelectLocation(avaliableLocations)
	refreshWeatherData := userSelectUpdate()
	if !refreshWeatherData {
		displayStaticWeather(selectedLocation)
		os.Exit(0)
	}

	pollWeatherDataLoop(selectedLocation)

}

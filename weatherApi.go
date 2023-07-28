package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type OpenMetoSearchNameResult struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Latitude    float32 `json:"latitude"`
	Longitude   float32 `json:"longitude"`
	Elevation   float32 `json:"elevation"`
	FeatureCode string  `json:"feature_code"`
	CounryCode  string  `json:"country_code"`
	Timezone    string  `json:"timezone"`
	Population  int     `json:"population"`
	CountryId   int     `json:"country_id"`
	Country     string  `json:"country"`
}

type OpenMetoSearchNameResponse struct {
	Results          []OpenMetoSearchNameResult `json:"results"`
	GenerationtimeMs float32                    `json:"generationtime_ms"`
	Error            bool                       `json:"error"`
	Reason           string                     `json:"reason"`
}

type OpenMetoCurrentWeatherResponse struct {
	Latitude         float32                      `json:"latitude"`
	Longitude        float32                      `json:"longitude"`
	GenerationtimeMs float32                      `json:"generationtime_ms"`
	UTCOffsetSeconds int                          `json:"utc_offset_seconds"`
	Timezone         string                       `json:"timezone"`
	Elevation        float32                      `json:"elevation"`
	Results          OpenMetoCurrentWeatherResult `json:"current_weather"`
	Error            bool                         `json:"error"`
	Reason           string                       `json:"reason"`
}

type OpenMetoCurrentWeatherResult struct {
	Temperature   float32 `json:"temperature"`
	Windspeed     float32 `json:"windspeed"`
	Winddirection float32 `json:"winddirection"`
	Weathercode   int     `json:"weathercode"`
	IsDay         int     `json:"is_day"`
	Time          string  `json:"time"`
}

type LocationCoordinates struct {
	Latitude  float32
	Longitude float32
}

const BASE_METO_GEOCODING_URL = "https://geocoding-api.open-meteo.com/v1/"
const BASE_METO_URL = "https://api.open-meteo.com/v1/"

func QuerySearchLocations(location string) ([]OpenMetoSearchNameResult, error) {
	searchUrl := BASE_METO_GEOCODING_URL + "search" + "?name=" + url.QueryEscape(location)
	res, err := http.Get(searchUrl)
	if err != nil {
		return []OpenMetoSearchNameResult{}, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body) // The is array of bytes
	var jsonData OpenMetoSearchNameResponse
	if err := json.Unmarshal(body, &jsonData); // Parse byte to go struct
	err != nil {
		return []OpenMetoSearchNameResult{}, err
	}

	if jsonData.Error == true {
		return nil, errors.New(fmt.Sprintf("API error reason %v", jsonData.Reason))
	}

	return jsonData.Results, nil
}

func QueryCurrentWeather(coordinates LocationCoordinates) (OpenMetoCurrentWeatherResult, error) {
	queryLattitude := strconv.FormatFloat(float64(coordinates.Latitude), 'f', -1, 32)
	queryLongitude := strconv.FormatFloat(float64(coordinates.Longitude), 'f', -1, 32)
	currentWeatherUrl := BASE_METO_URL + "forecast" +
		"?latitude=" +
		url.QueryEscape(queryLattitude) +
		"&longitude=" + url.QueryEscape(queryLongitude) +
		"&current_weather=true"

	res, err := http.Get(currentWeatherUrl)
	if err != nil {
		return OpenMetoCurrentWeatherResult{}, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body) // The is array of bytes
	var jsonData OpenMetoCurrentWeatherResponse
	if err := json.Unmarshal(body, &jsonData); err != nil {
		return OpenMetoCurrentWeatherResult{}, err
	}

	if jsonData.Error == true {
		return OpenMetoCurrentWeatherResult{}, errors.New(fmt.Sprintf("API error reason %v", jsonData.Reason))
	}

	return jsonData.Results, nil
}

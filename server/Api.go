package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/arvindpadev/weatherman/external"
)

const OwmEndpoint = "https://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&appid=%s"

func outsideWeatherCondition(lat string, lon string, apiKey string) (*WeatherCondition, int, error) {
	latitude, errLat := strconv.ParseFloat(lat, 64)
	if errLat != nil || latitude < -90 || latitude > 90 {
		return nil, http.StatusBadRequest, fmt.Errorf("Incorrect latitude %s. Latitudes must be between [-90, 90]", lat)
	}

	longitude, errLong := strconv.ParseFloat(lon, 64)
	if errLong != nil || longitude < -180 || longitude >= 180 {
		return nil, http.StatusBadGateway, fmt.Errorf("Incorrect longitude %s. Longitudes must be between [-180, 180)", lon)
	}

	url := fmt.Sprintf(OwmEndpoint, lat, lon, apiKey)
	external.Debug.Println(fmt.Sprintf("Calling external API %s", url))
	response, err := http.Get(url)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	defer response.Body.Close()
	r, err := io.ReadAll(response.Body)
	if err != nil {
		external.Debug.Println(fmt.Sprintf("API error %v", err))
		return nil, http.StatusInternalServerError, err
	}

	external.Debug.Println(fmt.Sprintf("API Raw Response %v", string(r)))
	if response.StatusCode/500 == 1 || response.StatusCode/400 == 1 || response.StatusCode/300 == 1 {
		return nil, response.StatusCode, fmt.Errorf(string(r))
	}

	if response.StatusCode != http.StatusOK {
		return nil, http.StatusInternalServerError, fmt.Errorf("Unexpected response %v %s", response.StatusCode, string(r))
	}

	var owmResponse OwmResponse
	errResponse := json.Unmarshal(r, &owmResponse)
	if errResponse != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Poorly formed response received %s", string(r))
	}

	if len(owmResponse.Weather) == 0 {
		return nil, http.StatusInternalServerError, fmt.Errorf("Unexpected response format received %s", string(r))
	}

	var feeling FeelsLike = Hot
	if owmResponse.Main.FeelsLike > 273.15 && owmResponse.Main.FeelsLike < 299.817 {
		feeling = Moderate
	} else if owmResponse.Main.FeelsLike <= 273.15 {
		feeling = Cold
	}

	return &WeatherCondition{
		Feeling:   feeling,
		Condition: owmResponse.Weather[0].Main,
	}, http.StatusOK, nil
}

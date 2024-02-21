package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/arvindpadev/weatherman/external"
)

func Run() {
	external.InitCustomLoggers()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /conditions", restOutsideWeatherConditions)
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func restOutsideWeatherConditions(w http.ResponseWriter, r *http.Request) {
	external.Debug.Println("intercepted")
	if !r.URL.Query().Has("lat") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Latitude not supplied"))
		return
	}

	if !r.URL.Query().Has("lon") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Longitude not supplied"))
		return
	}

	if !r.URL.Query().Has("appid") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("API Key not supplied"))
		return
	}

	lat := r.URL.Query().Get("lat")
	lon := r.URL.Query().Get("lon")
	apiKey := r.URL.Query().Get("appid")
	external.Debug.Println("Calling api function")
	conditions, statusCode, err := outsideWeatherCondition(lat, lon, apiKey)
	if err != nil {
		w.WriteHeader(statusCode)
		w.Write([]byte(err.Error()))
	} else {
		response, err := json.Marshal(conditions)
		if err != nil {
			external.Error.Println(fmt.Sprintf("For request %s, unable to marshal response for %v", r.RequestURI, *conditions))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Unable to marshal response"))
		} else {
			w.WriteHeader(statusCode)
			w.Write(response)
		}
	}
}

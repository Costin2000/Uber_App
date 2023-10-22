package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (app *Config) GetCarRequests(w http.ResponseWriter, r *http.Request) {
	bearer := r.Header.Get("Authorization")

	request, err := http.NewRequest("POST", "http://authentication-service/check_token", nil)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	if len(bearer) > 0 {
		request.Header.Set("Authorization", bearer)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		var payload errorResponse
		err = json.NewDecoder(response.Body).Decode(&payload)
		payload.Message = "Invalid token"
		app.writeJSON(w, response.StatusCode, payload)
		return
	}

	carType := r.URL.Query().Get("car_type")
	city := r.URL.Query().Get("city")
	activeStr := r.URL.Query().Get("active")
	active, err := strconv.ParseBool(activeStr)
	if err != nil {
		active = true
	}
	url := fmt.Sprintf("http://car-service/car_request?car_type=%s&city=%s&active=%t",
		carType, city, active)
	request, err = http.NewRequest("GET", url, nil)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	response, err = client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		var payload errorResponse
		err = json.NewDecoder(response.Body).Decode(&payload)
		app.writeJSON(w, response.StatusCode, payload)
		return
	}

	var jsonFromService jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Car requests retrieved"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)
}

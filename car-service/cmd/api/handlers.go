package main

import (
	"car-service/data"
	"fmt"
	"net/http"
	"strconv"
)

func (app *Config) CreateCar(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		UserId  int    `json:"user_id"`
		CarName string `json:"car_name"`
		City    string `json:"city"`
		CarType string `json:"car_type"`
	}

	logRequestBody(r)

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	car := data.Car{
		UserId:  requestPayload.UserId,
		CarName: requestPayload.CarName,
		City:    requestPayload.City,
		CarType: requestPayload.CarType,
	}

	_, err = app.Models.Car.InsertCar(car)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Car has been crated"),
		Data:    car,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) CreateCarRequest(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		UserId   int    `json:"user_id"`
		UserName string `json:"user_name"`
		CarType  string `json:"car_type"`
		City     string `json:"city"`
		Address  string `json:"address"`
	}

	logRequestBody(r)

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	carRequest := data.CarRequest{
		UserId:   requestPayload.UserId,
		UserName: requestPayload.UserName,
		City:     requestPayload.City,
		CarType:  requestPayload.CarType,
		Address:  requestPayload.Address,
	}

	_, err = app.Models.CarRequest.InsertCarRequest(carRequest)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Car has been crated"),
		Data:    carRequest,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) GetAllCarRequests(w http.ResponseWriter, r *http.Request) {
	carType := r.URL.Query().Get("car_type")
	city := r.URL.Query().Get("city")
	activeStr := r.URL.Query().Get("active")
	active, err := strconv.ParseBool(activeStr)
	if err != nil {
		active = true
	}

	carRequests, err := app.Models.CarRequest.GetAllCarRequestByCity(city, carType, active)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	type CarRequestsResponse struct {
		CarRequests []data.CarRequest `json:"car_requests"`
	}

	// Create a new slice of CarRequest with the same length as carRequests
	convertedCarRequests := make([]data.CarRequest, len(carRequests))

	// Copy values from carRequests (of type []*CarRequest) to convertedCarRequests
	for i, cr := range carRequests {
		convertedCarRequests[i] = *cr
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Car has been crated"),
		Data:    CarRequestsResponse{CarRequests: convertedCarRequests},
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

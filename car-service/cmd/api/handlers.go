package main

import (
	"car-service/data"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (app *Config) CreateCar(w http.ResponseWriter, r *http.Request) {
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
		app.errorJSON(w, errors.New("internal server error"))
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("invalid token"))
		return
	}

	var jsonFromServiceAuth jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromServiceAuth)

	tkData := jsonFromServiceAuth.Data.(map[string]interface{})
	userId := int(tkData["user_id"].(float64))

	userType := tkData["type"].(string)
	if userType != "driver" {
		app.errorJSON(w, errors.New("You are on a customer account. Should be logged in on a driver account to create cars."))
		return
	}

	var requestPayload struct {
		UserId  int    `json:"user_id"`
		CarName string `json:"car_name"`
		City    string `json:"city"`
		CarType string `json:"car_type"`
	}

	//logRequestBody(r)

	err = app.readJSON(w, r, &requestPayload)
	requestPayload.UserId = userId
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
		app.errorJSON(w, errors.New("internal server error"))
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("invalid token"))
		return
	}

	var jsonFromServiceAuth jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromServiceAuth)

	tkData := jsonFromServiceAuth.Data.(map[string]interface{})
	userId := int(tkData["user_id"].(float64))
	userName := tkData["username"].(string)

	var requestPayload struct {
		UserId   int    `json:"user_id"`
		UserName string `json:"user_name"`
		CarType  string `json:"car_type"`
		City     string `json:"city"`
		Address  string `json:"address"`
	}

	err = app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	requestPayload.UserId = userId
	requestPayload.UserName = userName

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
		Message: fmt.Sprintf("Car request has been crated"),
		Data:    carRequest,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) GetAllCarRequests(w http.ResponseWriter, r *http.Request) {
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
		app.errorJSON(w, errors.New("internal server error"))
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("invalid token"))
		return
	}

	var jsonFromServiceAuth jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromServiceAuth)

	tkData := jsonFromServiceAuth.Data.(map[string]interface{})
	userType := tkData["type"].(string)
	if userType != "driver" {
		app.errorJSON(w, errors.New("you are on a customer account. Should be logged in on a driver account to get the car requests"))
		return
	}

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

func (app *Config) GetAllCars(w http.ResponseWriter, r *http.Request) {
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
		app.errorJSON(w, errors.New("internal server error"))
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("invalid token"))
		return
	}

	var jsonFromServiceAuth jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromServiceAuth)

	tkData := jsonFromServiceAuth.Data.(map[string]interface{})
	userType := tkData["type"].(string)
	userId := int(tkData["user_id"].(float64))
	if userType != "driver" {
		app.errorJSON(w, errors.New("you are on a customer account. Should be logged in on a driver account to retrieve your cars"))
		return
	}

	activeStr := r.URL.Query().Get("active")
	active, err := strconv.ParseBool(activeStr)
	if err != nil {
		active = true
	}
	cars, err := app.Models.Car.GetAllCars(userId, active)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	type CarsResponse struct {
		Cars []data.Car `json:"cars"`
	}

	// Create a new slice of CarRequest with the same length as carRequests
	convertedCars := make([]data.Car, len(cars))

	// Copy values from carRequests (of type []*CarRequest) to convertedCarRequests
	for i, cr := range cars {
		convertedCars[i] = *cr
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Cars have been retrieved"),
		Data:    CarsResponse{Cars: convertedCars},
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) UpdateCar(w http.ResponseWriter, r *http.Request) {
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
		app.errorJSON(w, errors.New("internal server error"))
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("invalid token"))
		return
	}

	var jsonFromServiceAuth jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromServiceAuth)

	tkData := jsonFromServiceAuth.Data.(map[string]interface{})
	userType := tkData["type"].(string)
	if userType != "driver" {
		app.errorJSON(w, errors.New("you are on a customer account. Should be logged in on a driver account to update cars requests"))
		return
	}

	userId := int(tkData["user_id"].(float64))

	var requestPayload struct {
		Active bool `json:"active"`
	}

	err = app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	carId := chi.URLParam(r, "id")
	intCarId, _ := strconv.Atoi(carId)
	car, err := app.Models.Car.GetCarByID(intCarId)
	if err != nil {
		app.errorJSON(w, errors.New("car not found"), http.StatusBadRequest)
		return
	}

	if userId != car.UserId {
		app.errorJSON(w, errors.New("The car does not belong to you"), http.StatusBadRequest)
		return
	}

	car.Active = requestPayload.Active

	err = car.Update()
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("The car has been updated"),
		Data:    car,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) UpdateCarRequest(w http.ResponseWriter, r *http.Request) {
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
		app.errorJSON(w, errors.New("internal server error"))
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("invalid token"))
		return
	}

	var jsonFromServiceAuth jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromServiceAuth)

	tkData := jsonFromServiceAuth.Data.(map[string]interface{})
	userId := int(tkData["user_id"].(float64))
	userType := tkData["type"].(string)

	var requestPayload struct {
		Active bool `json:"active"`
		CarId  *int `json:"car_id"`
		Rating int  `json:"rating,omitempty"`
	}

	err = app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	carRequestId := chi.URLParam(r, "id")
	intCarRequestId, _ := strconv.Atoi(carRequestId)
	carRequest, err := app.Models.CarRequest.GetCarRequestByID(intCarRequestId)
	if err != nil {
		app.errorJSON(w, errors.New("car request not found"), http.StatusBadRequest)
		return
	}

	if userId != carRequest.UserId && userType != "driver" {
		app.errorJSON(w, errors.New("the car request does not belong to you and you are not a driver"), http.StatusBadRequest)
		return
	}

	carRequest.Active = requestPayload.Active
	if requestPayload.CarId == nil {
		carRequest.CarId = sql.NullInt64{Valid: false}
	} else {
		carRequest.CarId = sql.NullInt64{Int64: int64(*requestPayload.CarId), Valid: true}
	}
	carRequest.Rating = requestPayload.Rating

	err = carRequest.Update()
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("The car request has been updated"),
		Data:    carRequest,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

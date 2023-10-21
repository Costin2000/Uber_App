package main

import (
	"car-service/data"
	"fmt"
	"net/http"
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

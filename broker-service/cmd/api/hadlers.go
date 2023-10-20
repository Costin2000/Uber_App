package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action     string            `json:"action"`
	Auth       AuthPayload       `json:"auth,omitempty"`
	Register   RegisterPayload   `json:"register,omitempty"`
	UpdateUser UpdateUserPayload `json:"update_user,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterPayload struct {
	FirstName            string `json:"first_name,omitempty"`
	LastName             string `json:"last_name,omitempty"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
	City                 string `json:"city"`
	Type                 string `json:"type"`
}

type UpdateUserPayload struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email"`
	City      string `json:"city"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	bearer := r.Header.Get("Authorization")
	if len(bearer) > 0 {
		bearer = bearer[len("Bearer "):]
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "register":
		app.register(w, requestPayload.Register)
	case "edit_user":
		app.updateUser(w, requestPayload.UpdateUser, bearer)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create some json we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	//make sure we get back the correct status code
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
	payload.Message = "Authenticated"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) register(w http.ResponseWriter, a RegisterPayload) {
	// create some json we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the service
	request, err := http.NewRequest("POST", "http://authentication-service/register", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
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
		app.writeJSON(w, response.StatusCode, payload)
		return
	}

	var jsonFromService jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Registered"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) updateUser(w http.ResponseWriter, a UpdateUserPayload, bearer string) {
	// create some json we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the service
	request, err := http.NewRequest("PUT", "http://authentication-service/users", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	if len(bearer) > 0 {
		request.Header.Set("Authorization", "Bearer "+bearer)
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
		app.writeJSON(w, response.StatusCode, payload)
		return
	}

	var jsonFromService jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)

	var payload jsonResponse
	payload.Error = false
	payload.Message = "User updated successfully"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)
}

package main

import (
	"authentification/data"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	//validate the user agains the database
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("Invalid Credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("Invalid Credentials"), http.StatusBadRequest)
		return
	}

	type response struct {
		User  *data.User `json:"user"`
		Token string     `json:"token"`
	}

	tokenString, _ := createToken(user.FirstName+user.LastName, user.Email, user.ID, user.Type)

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data: response{
			User:  user,
			Token: tokenString,
		},
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) Register(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		FirstName            string `json:"first_name,omitempty"`
		LastName             string `json:"last_name,omitempty"`
		Email                string `json:"email"`
		City                 string `json:"city"`
		Type                 string `json:"type"`
		Password             string `json:"password"`
		PasswordConfirmation string `json:"password_confirmation"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	//check if email is uniq
	_, err = app.Models.User.GetByEmail(requestPayload.Email)
	if err == nil {
		// An error occurred during the database query
		app.errorJSON(w, errors.New("email already exists"), http.StatusBadRequest)
		return
	}

	if requestPayload.PasswordConfirmation != requestPayload.Password {
		app.errorJSON(w, errors.New("password and password confirmation do not match"), http.StatusBadRequest)
		return
	}

	if len(requestPayload.LastName) == 0 || len(requestPayload.FirstName) == 0 {
		app.errorJSON(w, errors.New("first name and last name should not be empty"), http.StatusBadRequest)
		return
	}

	if len(requestPayload.Password) < 6 {
		app.errorJSON(w, errors.New("password should be at least 6 characters long"), http.StatusBadRequest)
		return
	}

	newUser := data.User{
		Email:     requestPayload.Email,
		Password:  requestPayload.Password,
		FirstName: requestPayload.FirstName,
		LastName:  requestPayload.LastName,
		City:      requestPayload.City,
		Type:      requestPayload.Type,
	}

	_, err = app.Models.User.Insert(newUser)

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprint("User registered successfully"),
		Data:    newUser,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) Update(w http.ResponseWriter, r *http.Request) {
	for key, values := range r.Header {
		for _, value := range values {
			// Print each header key and its values
			log.Printf("Header: %s = %s", key, value)
		}
	}

	var requestPayload struct {
		Email     string `json:"email"`
		FirstName string `json:"first_name,omitempty"`
		LastName  string `json:"last_name,omitempty"`
		City      string `json:"city"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		app.errorJSON(w, errors.New("Missing authorization header"), http.StatusUnauthorized)
		return
	}
	tokenString = tokenString[len("Bearer "):]

	err = verifyToken(tokenString)
	if err != nil {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	err = app.checkTokenData(tokenString)
	if err != nil {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("user not found"), http.StatusNotFound)
		return
	}

	if len(requestPayload.LastName) == 0 || len(requestPayload.FirstName) == 0 {
		app.errorJSON(w, errors.New("first name and last name should not be empty"), http.StatusBadRequest)
		return
	}

	user.City = requestPayload.City
	user.FirstName = requestPayload.FirstName
	user.LastName = requestPayload.LastName

	err = user.Update()
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprint("User registered successfully"),
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

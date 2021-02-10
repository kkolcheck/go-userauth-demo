// File: login-handlers.go
package main

import (
	"encoding/json"
	"net/http"
	"time"
	"strconv"
	"fmt"
	"os"
	"io/ioutil"
	"errors"
	"github.com/jinzhu/copier"
)

// Payload represents the expected payload from the user.
type Payload struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token int `json:"token"`
}

// Users represents a collection of users.
type Users struct {
	Users []User `json:"users"`
}

// User represents the valid user login credentials in our system.
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Return an array of users.
func getUsers() (users Users, err error) {
	jsonFile, err := os.Open("json/stub-user-credentials.json")
	if err != nil {
		err = errors.New("Error opening user file")
		return
	}
	defer jsonFile.Close()

	byteValue, readErr := ioutil.ReadAll(jsonFile)
	if readErr != nil {
		err = errors.New("Error reading user file")
		return
	}

	unmarshalErr := json.Unmarshal(byteValue, &users)
	if unmarshalErr != nil {
		err = errors.New("Error unmarsheling json")
		return
	}
	return
}

// Return user if user is found among the list of users.
func getUser(p Payload) (User, error){
	// If you are using a db to do the look up, this match would be a select query.
	var user User
	users, loadErr := getUsers()
	if loadErr != nil {
		fmt.Printf("Unable to load users: %s\n", loadErr.Error())
		return user, loadErr
	}

	for _, u := range users.Users {
		if (p.Username == u.Username && p.Password == u.Password) {
			copier.Copy(&user, &u)
			return user, nil
		}
	}
	fmt.Println("User not found")
	notFoundErr := errors.New("User not found")
	return user, notFoundErr
}

// Return boolean indicating that the token matches the current time in hhmm format. Leading zeroes are dropped.
func isValidToken(p Payload) bool {
	// Set timezone to ensure generated and validation tokens match.
	loc, _ := time.LoadLocation("America/New_York")
	now := time.Now().In(loc)

	h := strconv.Itoa(now.Hour())
	m := strconv.Itoa(now.Minute())
	if (now.Minute() < 10) {
		m = "0" + m
	}
	strToken := h + m
	token, _ := strconv.Atoi(strToken)
	if p.Token != token {
		fmt.Println(`Invalid Token`)
		return false;
	}
	return true;
}

// Return an error if payload decode fails.
func decodePayload(p *Payload, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		fmt.Println(`Invalid JSON`)
		return err
	}
	return nil	
}

// POST Login Handler will validate a login request.
// The handler expects a payload containing a user's username, password, and a valid security token.
func postLoginHandler(w http.ResponseWriter, r *http.Request) {
	var p Payload
	decodeErr := decodePayload(&p, r)
	if (decodeErr != nil || !isValidToken(p)) {
		http.Error(w, "{ \"message\": \"Bad Request\" }", http.StatusBadRequest)
		return
	}

	_, err := getUser(p)
	if err != nil {
		if (err.Error() == "Error opening user file" ||
			err.Error() == "Error reading user file" ||
			err.Error() == "Error unmarsheling json") {
			http.Error(w, "{ \"message\": \"Internal Server Error\" }", http.StatusInternalServerError)
			return 
		}
		http.Error(w, "{ \"message\": \"Not Found\" }", http.StatusNotFound)
		return 
	}

	fmt.Println(`User authenticated`)
	w.Write([]byte("{ \"message\": \"Success\" }"))
}

// Route request to appropriate handler.
func loginHandler(w http.ResponseWriter, r *http.Request) {
	setupResponse(w)
	switch r.Method {
	case "OPTIONS":
		return
	case "POST":
		postLoginHandler(w, r)
	default:
		http.Error(w, "{ \"message\": \"Not Implemented\" }", http.StatusNotImplemented)
	}
}

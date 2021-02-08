// File: login-handlers.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"strconv"
	"fmt"
	"os"
	"io/ioutil"
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

func getUsers() Users {
	// using a stub in place of a db
	jsonFile, err := os.Open("json/stub-user-credentials.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var users Users
	json.Unmarshal(byteValue, &users)
	return users
}

func isValidUser(p *Payload) bool {
	users := getUsers()

	for _, u := range users.Users {
		if (p.Username == u.Username && p.Password == u.Password) {
			return true
		}
	}
	return false
}

func isValidToken(p *Payload) bool {
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
		return false;
	}
	return true;
}

func isValidPayload(w *http.ResponseWriter, r *http.Request) bool {
	decoder := json.NewDecoder(r.Body)

	var p Payload
	err := decoder.Decode(&p)

	if (err != nil) {
		fmt.Println(`Invalid JSON`)
		http.Error(*w, `"{ \"message\": \"Bad Request\" }"`, http.StatusBadRequest)
		return false
	}
	if (!isValidToken(&p)) {
		fmt.Println(`Bad Request`)
		http.Error(*w, "{ \"message\": \"Bad Requeste\" }", http.StatusBadRequest)
		return false
	}

	if (!isValidUser(&p)) {
		fmt.Println(`Not Found`)
		http.Error(*w, "{ \"message\": \"Not Found\" }", http.StatusNotFound)
		return false
	}
	return true	
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w)
	switch r.Method {
	case "OPTIONS":
		return
	case "POST":
		if (!isValidPayload(&w, r)) {
			return
		}
		w.Write([]byte("{ \"message\": \"Success\" }"))
	default:
		http.Error(w, "{ \"message\": \"Not Implemented\" }", http.StatusNotImplemented)
	}
}
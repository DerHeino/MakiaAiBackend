package main

import (
	"encoding/json"
	"errors"
	"health/route"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"regexp"
	"strings"
)

// This file contains functions for converting http.Requests into a format
// for further processing, or auth validation

// Converts supported content types of HTTP request into a map for further processing
// Only required for POST requests
//
// returns a map with information and empty error if successful
func ConvertRequest(r *http.Request) (map[string]interface{}, error) {

	contentType := r.Header.Get("Content-Type")
	if match, _ := regexp.MatchString("application/json", contentType); match {
		return JsonToMap(r)
	} else if match, _ := regexp.MatchString("multipart/form-data", contentType); match {
		return MultiFormToMap(r)
	} else if match, _ := regexp.MatchString("application/x-www-form-urlencoded", contentType); match {
		return FormToMap(r)
	}

	return nil, errors.New("Unsupported Content-Type: " + contentType)
}

func JsonToMap(r *http.Request) (map[string]interface{}, error) {
	requestMap := make(map[string]interface{})

	err := json.NewDecoder(r.Body).Decode(&requestMap)
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("JSON formatting error")
	}
	return requestMap, nil
}

func FormToMap(r *http.Request) (map[string]interface{}, error) {
	requestMap := make(map[string]interface{})

	err := r.ParseForm()
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("Form formatting error")
	}

	for k, v := range r.Form {
		requestMap[k] = v[0]
	}

	return requestMap, nil
}

func MultiFormToMap(r *http.Request) (map[string]interface{}, error) {
	requestMap := make(map[string]interface{})

	err := r.ParseMultipartForm(r.ContentLength)
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("Form formatting error")
	}

	for k, v := range r.Form {
		requestMap[k] = v[0]
	}

	return requestMap, nil
}

// For image processing
func MultiFormImage(r *http.Request) (image.Image, error) {

	err := r.ParseMultipartForm(r.ContentLength)
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("Form formatting error")
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		return nil, errors.New("Form formatting error")
	}

	// image in memory
	deviceImage, err := jpeg.Decode(file)
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("failed to decode jpg image")
	}

	// Testing creates and writes into root of module
	//f, _ := os.Create("testX.jpg")
	//defer f.Close()
	//if err = jpeg.Encode(f, deviceImage, nil); err != nil {
	//	log.Printf("failed to encode: %v", err)
	//}

	//fmt.Println(reflect.TypeOf(h), reflect.TypeOf(file))
	return deviceImage, nil
}

// Takes "Authorization" string from the HTTP header and
// evaluates its parsing ("Bearer" <Token>)
// Takes into account if admin rights are required for POST /project and /location
// calls ValidateToken from route/user.go
//
// returns a string and an error for ResponseWriter and error logging
// if the token is invalid in any way
// returns an empty string and error if token is valid
func ValidateTokenForm(token string, admin bool) error {

	parts := strings.Split(token, "Bearer")
	if len(parts) != 2 {
		return errors.New("no Token found")
	} else {
		tokenString := strings.TrimSpace(parts[1])
		//fmt.Println(token)

		if ok, _ := route.ValidateToken(tokenString); ok {
			return nil
		} else {
			return errors.New("access denied, Token invalid: " + tokenString)
		}
	}
}

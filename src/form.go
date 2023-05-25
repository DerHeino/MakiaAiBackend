package main

import (
	"encoding/json"
	"errors"
	c "health/clog"
	"image"
	"image/jpeg"
	"net/http"
	"regexp"
)

// Converts supported content types of HTTP request into a map for further processing
// Only required for POST requests
//
// returns a map with data and empty error if successful
func ConvertRequest(r *http.Request) (map[string]interface{}, error) {

	contentType := r.Header.Get("Content-Type")
	if match, _ := regexp.MatchString("application/json", contentType); match {
		return JsonToMap(r)
	} else if match, _ := regexp.MatchString("multipart/form-data", contentType); match {
		return MultiFormToMap(r)
	} else if match, _ := regexp.MatchString("application/x-www-form-urlencoded", contentType); match {
		return FormToMap(r)
	}

	return nil, errors.New("unsupported Content-Type: " + contentType)
}

func JsonToMap(r *http.Request) (map[string]interface{}, error) {
	requestMap := make(map[string]interface{}, 20)

	err := json.NewDecoder(r.Body).Decode(&requestMap)
	if err != nil {
		c.ErrorLog.Println(err.Error())
		return nil, errors.New("JSON formatting error")
	}
	return requestMap, nil
}

func FormToMap(r *http.Request) (map[string]interface{}, error) {
	requestMap := make(map[string]interface{}, 20)

	err := r.ParseForm()
	if err != nil {
		c.ErrorLog.Println(err.Error())
		return nil, errors.New("form formatting error")
	}

	for k, v := range r.Form {
		requestMap[k] = v[0]
	}

	return requestMap, nil
}

func MultiFormToMap(r *http.Request) (map[string]interface{}, error) {
	requestMap := make(map[string]interface{}, 20)

	err := r.ParseMultipartForm(r.ContentLength)
	if err != nil {
		c.ErrorLog.Println(err.Error())
		return nil, errors.New("form formatting error")
	}

	for k, v := range r.Form {
		requestMap[k] = v[0]
	}

	return requestMap, nil
}

// Image uploads support only one content type as well as only one image format (jpeg)
func MultiFormImage(r *http.Request) (image.Image, error) {

	err := r.ParseMultipartForm(r.ContentLength)
	if err != nil {
		c.ErrorLog.Println(err.Error())
		return nil, errors.New("form formatting error")
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		c.ErrorLog.Println(err.Error())
		return nil, errors.New("form formatting error")
	}

	// image in memory
	deviceImage, err := jpeg.Decode(file)
	if err != nil {
		c.ErrorLog.Println(err.Error())
		return nil, errors.New("failed to decode jpg image")
	}

	return deviceImage, nil
}

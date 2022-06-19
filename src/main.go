package main

import (
	"bytes"
	"fmt"
	"health/network"
	"health/route"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/go-chi/chi/v5"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HealthCheck-API \"/\"")
	fmt.Println("Endpoint Hit: home")
}

func login(w http.ResponseWriter, r *http.Request) {

	mapRequest, err := ConvertRequest(r)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	tokenString, err := route.VerifyLogin(mapRequest)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	fmt.Fprint(w, tokenString)
}

func getProject(w http.ResponseWriter, r *http.Request) {
	if err := ValidateTokenForm(r.Header.Get("Authorization"), false); err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	projectsJson, err := network.GetAllDocuments("projects")
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	fmt.Fprint(w, string(projectsJson))
}

func postProject(w http.ResponseWriter, r *http.Request) {
	if err := ValidateTokenForm(r.Header.Get("Authorization"), true); err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	mapRequest, err := ConvertRequest(r)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	confirmation, err := route.PostProject(mapRequest)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	fmt.Fprintf(w, "project %s", confirmation)
}

func getLocation(w http.ResponseWriter, r *http.Request) {
	if err := ValidateTokenForm(r.Header.Get("Authorization"), false); err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	projectsJson, err := network.GetAllDocuments("locations")
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	fmt.Fprint(w, string(projectsJson))
}

func postLocation(w http.ResponseWriter, r *http.Request) {
	if err := ValidateTokenForm(r.Header.Get("Authorization"), true); err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	mapRequest, err := ConvertRequest(r)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	confirmation, err := route.PostLocation(mapRequest)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	fmt.Fprintf(w, "location %s", confirmation)
}

func getDevice(w http.ResponseWriter, r *http.Request) {
	if err := ValidateTokenForm(r.Header.Get("Authorization"), false); err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	projectsJson, err := network.GetAllDocuments("devices")
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	fmt.Fprint(w, string(projectsJson))
}

func postDevice(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.Header)
	if err := ValidateTokenForm(r.Header.Get("Authorization"), false); err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	mapRequest, err := ConvertRequest(r)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	confirmation, err := route.PostDevice(mapRequest)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	fmt.Fprintf(w, "device %s", confirmation)
}

func getInventory(w http.ResponseWriter, r *http.Request) {
	if err := ValidateTokenForm(r.Header.Get("Authorization"), false); err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	projectsJson, err := network.GetAllDocuments("inventory")
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	fmt.Fprint(w, string(projectsJson))
}

func postInventory(w http.ResponseWriter, r *http.Request) {
	if err := ValidateTokenForm(r.Header.Get("Authorization"), false); err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	mapRequest, err := ConvertRequest(r)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	confirmation, err := route.PostInventory(mapRequest)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	fmt.Fprintf(w, "inventory %s", confirmation)
}

func postPing(w http.ResponseWriter, r *http.Request) {

	pingRequest, err := ConvertRequest(r)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	confirmation, err := route.PostPing(pingRequest)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	fmt.Fprint(w, confirmation)
}

func getImage(w http.ResponseWriter, r *http.Request) {
	if err := ValidateTokenForm(r.Header.Get("Authorization"), false); err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	deviceId := chi.URLParam(r, "deviceId")

	//fmt.Println(network.DeviceList)

	image := route.GetImage(deviceId)
	if image == nil {
		fmt.Fprintf(w, "no image found under device %s", deviceId)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "image/jpg")

	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, *image, nil)
	if err != nil {
		fmt.Fprintf(w, "error displaying image")
		return
	}

	w.Write(buf.Bytes())

	//fmt.Fprintf(w, "deviceId %s", deviceId)
}

func uploadImage(w http.ResponseWriter, r *http.Request) {
	deviceId := chi.URLParam(r, "deviceId")

	contentType := r.Header.Get("Content-Type")

	if match, _ := regexp.MatchString("multipart/form-data", contentType); match {
		myImage, err := MultiFormImage(r)
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		if route.PostImage(deviceId, &myImage) {
			fmt.Fprintf(w, "image -> %s", deviceId)
		} else {
			fmt.Fprintf(w, "upload failed")
		}
	} else {
		fmt.Fprint(w, "Unsupported Content-Type: "+contentType)
	}
}

func handleRequests() {
	r := chi.NewRouter()

	r.Get("/", homePage)

	r.Post("/login", login)

	r.Get("/project", getProject)
	r.Post("/project", postProject)

	r.Get("/location", getLocation)
	r.Post("/location", postLocation)

	r.Get("/device", getDevice)
	r.Post("/device", postDevice)

	r.Get("/inventory", getInventory)
	r.Post("/inventory", postInventory)

	r.Post("/ping", postPing)
	r.Get("/device/{deviceId}/image", getImage)
	r.Post("/device/{deviceId}/image", uploadImage)

	port := ":" + os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(port, r))
}

func main() {
	setConfig()

	// build connection to firestore database
	client := network.Start_firebase()

	handleRequests()

	//_ = client
	defer client.Close()
}

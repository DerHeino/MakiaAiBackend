package main

import (
	"bytes"
	"fmt"
	bg "health/background"
	c "health/clog"
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
	fmt.Println("welcome")
}

func login(w http.ResponseWriter, r *http.Request) {

	mapRequest, err := ConvertRequest(r)
	if err != nil {
		fmt.Fprint(w, "error: ", err.Error())
		return
	}

	tokenString, err := route.VerifyLogin(mapRequest)
	if err != nil {
		fmt.Fprint(w, "error: ", err.Error())
		return
	}
	fmt.Fprint(w, tokenString)
}

func postRegister(w http.ResponseWriter, r *http.Request) {
	inviter, err := route.ValidateToken(r.Header.Get("Authorization"), true, os.Getenv("REGISTER_KEY"))
	if err != nil {
		fmt.Fprint(w, "error: ", err.Error())
		return
	}

	mapRequest, err := ConvertRequest(r)
	if err != nil {
		fmt.Fprint(w, "error: ", err.Error())
		return
	}

	tokenString, err := route.VerifyUser(inviter, mapRequest)
	if err != nil {
		fmt.Fprint(w, "error: ", err.Error())
		return
	}
	fmt.Fprint(w, tokenString)
}

func getRegisterKey(w http.ResponseWriter, r *http.Request) {
	username, err := route.ValidateToken(r.Header.Get("Authorization"), true)
	if err != nil {
		fmt.Fprint(w, "error: ", err.Error())
		return
	}

	if key, err := route.BuildRegisterKey(username); err != nil {
		fmt.Fprint(w, "error: ", err.Error())
	} else {
		fmt.Fprint(w, key)
	}
}

func getProject(w http.ResponseWriter, r *http.Request) {
	getModel(&w, r, "project")
}

func getLocation(w http.ResponseWriter, r *http.Request) {
	getModel(&w, r, "location")
}

func getDevice(w http.ResponseWriter, r *http.Request) {
	getModel(&w, r, "device")
}

func getInventory(w http.ResponseWriter, r *http.Request) {
	getModel(&w, r, "inventory")
}

func getModel(w *http.ResponseWriter, r *http.Request, model string) {
	if _, err := route.ValidateToken(r.Header.Get("Authorization"), false); err != nil {
		fmt.Fprint(*w, "error: ", err.Error())
		return
	}

	modelJson, err := network.GetAllDocuments(model)
	if err != nil {
		fmt.Fprint(*w, "error: ", err.Error())
		return
	}

	(*w).Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(*w, string(modelJson))
}

func postProject(w http.ResponseWriter, r *http.Request) {
	if _, err := route.ValidateToken(r.Header.Get("Authorization"), true); err != nil {
		fmt.Fprint(w, "error: ", err.Error())
		return
	}

	mapRequest, err := ConvertRequest(r)
	if err != nil {
		fmt.Fprint(w, "error: ", err.Error())
		return
	}

	projectId, err := route.PostProject(mapRequest)
	if err != nil {
		fmt.Fprint(w, "error: ", err.Error())
		return
	}
	fmt.Fprintf(w, "%s", projectId)
}

func postLocation(w http.ResponseWriter, r *http.Request) {
	if _, err := route.ValidateToken(r.Header.Get("Authorization"), true); err != nil {
		fmt.Fprint(w, "error: ", err.Error())
		return
	}

	mapRequest, err := ConvertRequest(r)
	if err != nil {
		fmt.Fprint(w, "error: ", err.Error())
		return
	}

	locId, err := route.PostLocation(mapRequest)
	if err != nil {
		fmt.Fprint(w, "error: ", err.Error())
		return
	}
	fmt.Fprintf(w, "%s", locId)
}

func postDevice(w http.ResponseWriter, r *http.Request) {
	if _, err := route.ValidateToken(r.Header.Get("Authorization"), false); err != nil {
		fmt.Fprint(w, "error: ", err.Error())
		return
	}

	mapRequest, err := ConvertRequest(r)
	if err != nil {
		fmt.Fprint(w, "error: ", err.Error())
		return
	}

	devId, err := route.PostDevice(mapRequest)
	if err != nil {
		if devId != "" {
			fmt.Fprint(w, devId, "\nerror: ", err.Error())
		} else {
			fmt.Fprint(w, "error: ", err.Error())
		}
		return
	}

	fmt.Fprintf(w, "%s", devId)
}

func postInventory(w http.ResponseWriter, r *http.Request) {
	if _, err := route.ValidateToken(r.Header.Get("Authorization"), false); err != nil {
		fmt.Fprint(w, "error: ", err.Error())
		return
	}

	mapRequest, err := ConvertRequest(r)
	if err != nil {
		fmt.Fprint(w, "error: ", err.Error())
		return
	}

	invId, err := route.PostInventory(mapRequest)
	if err != nil {
		fmt.Fprint(w, "error: ", err.Error())
		return
	}
	fmt.Fprintf(w, "%s", invId)
}

func postPing(w http.ResponseWriter, r *http.Request) {

	pingRequest, err := ConvertRequest(r)
	if err != nil {
		fmt.Fprint(w, "error: ", err.Error())
		return
	}

	lastPing, err := route.PostPing(pingRequest)
	if err != nil {
		fmt.Fprint(w, "error: ", err.Error())
		return
	}
	fmt.Fprint(w, lastPing)
}

func getImage(w http.ResponseWriter, r *http.Request) {
	/*if _, err := route.ValidateToken(r.Header.Get("Authorization"), false); err != nil {
		fmt.Fprint(w, "error: ", err.Error())
		return
	}*/

	deviceId := chi.URLParam(r, "deviceId")

	//fmt.Println(network.DeviceList)

	image := route.GetImage(deviceId)
	if image == nil {
		fmt.Fprintf(w, "error: no image found under device %s", deviceId)
		return
	}

	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, *image, nil)
	if err != nil {
		fmt.Fprintf(w, "error: displaying image")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "image/jpg")
	w.Write(buf.Bytes())
}

func uploadImage(w http.ResponseWriter, r *http.Request) {
	deviceId := chi.URLParam(r, "deviceId")

	contentType := r.Header.Get("Content-Type")

	if match, _ := regexp.MatchString("multipart/form-data", contentType); match {
		myImage, err := MultiFormImage(r)
		if err != nil {
			fmt.Fprint(w, "error: ", err.Error())
			return
		}
		if route.PostImage(deviceId, &myImage) {
			fmt.Fprintf(w, "image -> %s", deviceId)
		} else {
			fmt.Fprintf(w, "error: upload failed")
		}
	} else {
		fmt.Fprint(w, "error: unsupported Content-Type: ", contentType)
	}
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	if _, err := route.ValidateToken(r.Header.Get("Authorization"), true); err != nil {
		fmt.Fprint(w, "error: ", err.Error())
		return
	}

	id := chi.URLParam(r, "projectId")
	project := make([]byte, 0, 120)

	if err := route.DeleteProject(id, &project); err != nil {
		fmt.Fprint(w, "error: ", err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, string(project))
}

func deleteLocation(w http.ResponseWriter, r *http.Request) {
	if _, err := route.ValidateToken(r.Header.Get("Authorization"), true); err != nil {
		fmt.Fprint(w, "error: ", err.Error())
		return
	}

	id := chi.URLParam(r, "locationId")
	location := make([]byte, 0, 360)

	if err := route.DeleteLocation(id, &location, nil); err != nil {
		fmt.Fprint(w, "error: ", err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, string(location))
}

func deleteDevice(w http.ResponseWriter, r *http.Request) {
	if _, err := route.ValidateToken(r.Header.Get("Authorization"), false); err != nil {
		fmt.Fprint(w, "error: ", err.Error())
		return
	}

	id := chi.URLParam(r, "deviceId")
	device := make([]byte, 0, 360)

	if err := route.DeleteDevice(id, &device, nil); err != nil {
		fmt.Fprint(w, "error: ", err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, string(device))
}

func deleteInventory(w http.ResponseWriter, r *http.Request) {
	if _, err := route.ValidateToken(r.Header.Get("Authorization"), false); err != nil {
		fmt.Fprint(w, "error: ", err.Error())
		return
	}

	id := chi.URLParam(r, "inventoryId")
	inventory := make([]byte, 0, 240)

	if err := route.DeleteInventory(id, &inventory, nil); err != nil {
		//error message contains id, so test will always be successful
		fmt.Fprint(w, "error: ", err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, string(inventory))
}

func handleRequests() {
	r := chi.NewRouter()

	r.Get("/", homePage)

	r.Post("/login", login)
	r.Post("/register", postRegister)
	r.Get("/register", getRegisterKey)

	r.Get("/project", getProject)
	r.Get("/location", getLocation)
	r.Get("/device", getDevice)
	r.Get("/inventory", getInventory)

	r.Post("/project", postProject)
	r.Post("/location", postLocation)
	r.Post("/device", postDevice)
	r.Post("/inventory", postInventory)
	r.Post("/ping", postPing)

	r.Get("/device/{deviceId}/image", getImage)
	r.Post("/device/{deviceId}/image", uploadImage)

	r.Delete("/project/{projectId}", deleteProject)
	r.Delete("/location/{locationId}", deleteLocation)
	r.Delete("/device/{deviceId}", deleteDevice)
	r.Delete("/inventory/{inventoryId}", deleteInventory)

	port := ":" + os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(port, r))
}

func main() {
	setConfig()
	initLogs()

	c.InfoLog.Println("Starting Backend")

	// build connection to firestore database
	network.Start_firebase()

	// register map (removes entries after expiry)
	bg.RunRegMap()
	// image map
	bg.RunDevMap()

	handleRequests()
}

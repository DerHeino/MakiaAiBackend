package network

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"health/model"
	"log"
	"os"
	"time"

	"google.golang.org/api/option"

	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

var global_ctx context.Context      // initialized once
var global_client *firestore.Client // initialized once
var DeviceList map[string]model.LocalDevice

func Start_firebase() *firestore.Client {

	d, _ := json.Marshal(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	fmt.Println(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	// Use a service account
	ctx := context.Background()
	sa := option.WithCredentialsJSON(d)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	global_ctx = ctx
	global_client = client

	helloworld, err := client.Collection("_online").Doc("confirmation").Get(ctx)
	if err != nil {
		fmt.Println(err)
	}

	confirmation := helloworld.Data()
	fmt.Println(confirmation)

	CountDevices()

	return client
}

func CountDevices() {
	DeviceList = make(map[string]model.LocalDevice)

	deviceIterator := global_client.Collection("devices").Documents(global_ctx)

	for {
		device, err := deviceIterator.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return
		}
		id, err := device.DataAt("_id")
		if err != nil {
			log.Print(err.Error())
			return
		}
		if id, ok := id.(string); ok {
			DeviceList[id] = model.LocalDevice{}
		}
	}
}

// Retrieves user Doc from Firestore database and returns in
// struct format (alongside an error)
func GetUser(username string, retrieved *model.UserCredentials) error {

	userRef := global_client.Collection("users")
	credentials, err := userRef.Doc("0x" + fmt.Sprintf("%X", username)).Get(global_ctx)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	credentials.DataTo(retrieved)
	return nil
}

// for demonstration purposes
func setUser() {
	u := &model.UserCredentials{
		Username: "d.duck",
		Password: "$2a$10$nHxRig3l/6fsu8.fVf7eP.ycwR7xj8wgnTI7yskHJ8Pj4Crzlq3wO",
		User: &model.User{
			Id:    fmt.Sprintf("%X", "d.duck"),
			Name:  "Donald",
			Admin: false,
		},
	}

	_, err := global_client.Collection("users").Doc("0x"+u.User.Id).Set(global_ctx, *u)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// "/project", "/location", "/device", "/inventory" GET requests
// all return an array with its supported model.
//
// This method retrieves all documents from the specific collection
// converts them into an array[map] and returns it the JSON encoding of said array
func GetAllDocuments(collection string) ([]byte, error) {
	var array []map[string]interface{}

	docsIterator := global_client.Collection(collection).Documents(global_ctx)

	for {
		docs, err := docsIterator.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		array = append(array, docs.Data())
	}

	jsonStr, err := json.Marshal(array)
	if err != nil {
		return nil, err
	}
	return jsonStr, nil
}

func SetProjectFire(project *model.Project) bool {

	_, err := global_client.Collection("projects").Doc(project.Id).Set(global_ctx, *project)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	return true
}

func SetLocationFire(location *model.Location) bool {

	_, err := global_client.Collection("locations").Doc(location.Id).Set(global_ctx, *location)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	return true
}

func SetDeviceFire(device *model.Device) error {

	_, lerr := global_client.Collection("locations").Doc(device.LocationId).Get(global_ctx)
	if lerr != nil {
		log.Println(lerr.Error())
		if status.Code(lerr) == codes.NotFound {
			return errors.New("location " + device.LocationId + " was not found")
		}
		return lerr
	}

	_, err := global_client.Collection("devices").Doc(device.Id).Set(global_ctx, *device)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if _, ok := DeviceList[device.Id]; !ok {
		DeviceList[device.Id] = model.LocalDevice{}
	}

	return nil
}

func SetInventoryFire(inventory *model.Inventory) bool {

	_, err := global_client.Collection("inventory").Doc(inventory.Id).Set(global_ctx, *inventory)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

func UpdatePingFire(ping *model.Ping) (*model.Ping, error) {

	devSnap, _ := global_client.Collection("devices").Doc(ping.Id).Get(global_ctx)
	pingSnap, _ := devSnap.DataAt("lastPing")

	_, err := global_client.Collection("devices").Doc(ping.Id).Update(global_ctx, []firestore.Update{
		{
			Path:  "lastPing",
			Value: ping,
		},
	})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if val, ok := pingSnap.(map[string]interface{}); ok {
		ping.Timestamp = val["timestamp"].(time.Time)
	}

	return ping, nil
}

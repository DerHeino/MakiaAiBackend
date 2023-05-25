package network

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	c "health/clog"
	"health/model"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var global_ctx context.Context      // initialized once
var global_client *firestore.Client // initialized once
var DeviceList map[string]model.LocalImage
var AdminList = []string{}

func Start_firebase() *firestore.Client {

	//fmt.Println(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	// Use a service account
	ctx := context.Background()
	sa := option.WithCredentialsJSON([]byte(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")))
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
		c.ErrorLog.Println(err)
	}

	confirmation := helloworld.Data()
	fmt.Println(confirmation)

	CountDevices()
	setUser()

	return client
}

func CountDevices() {
	DeviceList = make(map[string]model.LocalImage)

	deviceIterator := global_client.Collection("device").Documents(global_ctx)

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
			c.ErrorLog.Println(err.Error())
			return
		}
		if id, ok := id.(string); ok {
			DeviceList[id] = model.LocalImage{}
		}
	}
}

func GetAllDeviceIDs() []string {
	var devices []string

	deviceIterator := global_client.Collection("device").Documents(global_ctx)

	for {
		device, err := deviceIterator.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil
		}

		if id, err := device.DataAt("_id"); err == nil {
			if id, ok := id.(string); ok {
				devices = append(devices, id)
				DeviceList[id] = model.LocalImage{}
			}
		} else {
			c.ErrorLog.Println(err.Error())
		}
	}

	return devices
}

func UserValid(userid string, admin bool) error {

	dsnap, err := global_client.Collection("user").Doc(userid).Get(global_ctx)
	if err != nil {
		c.ErrorLog.Println(err.Error())
		return fmt.Errorf("invalid token")
	}

	if admin {
		credentials := dsnap.Data()
		if !credentials["admin"].(bool) {
			return fmt.Errorf("access denied")
		}
	}

	return nil
}

// Retrieves user Doc from Firestore database and returns in
// struct format (alongside an error)
func GetUser(userid string, retrieved *model.UserCredentials) error {

	userRef := global_client.Collection("user")
	credentials, err := userRef.Doc(userid).Get(global_ctx)
	if err != nil {
		c.ErrorLog.Println(err.Error())
		return err
	}

	credentials.DataTo(retrieved)
	return nil
}

// for demonstration purposes
func setUser() {
	u := &model.UserCredentials{
		Username: "h.chan",
		Password: "$2a$10$G923.ZZPTWZ27YZHjWEJSOb/h7KSZQ0rdkckPIDmLgfXuSXVVNad6",
		User: model.User{
			Id:        "AC37846C9E8A4568CBDE218B",
			Name:      "Heinrich",
			Telephone: "+00/123456789",
			Admin:     true,
		},
	}

	_, err := global_client.Collection("user").Doc(u.User.Id).Set(global_ctx, *u)
	if err != nil {
		c.ErrorLog.Println(err.Error())
	}
}

// "/project", "/location", "/device", "/inventory" GET requests
// all return an array with its supported model.
//
// This method retrieves all documents from the specific collection
// converts them into an array[map] and returns the JSON encoding of said array
func GetAllDocuments(collection string) ([]byte, error) {
	var err error
	var array []map[string]interface{}

	docsIterator := global_client.Collection(collection).Documents(global_ctx)

	for {
		docs, err := docsIterator.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.ErrorLog.Println(err.Error())
			return nil, err
		}
		array = append(array, docs.Data())
	}

	jsonStr, err := json.Marshal(array)
	if err != nil {
		c.ErrorLog.Println(err.Error())
		return nil, err
	}
	return jsonStr, nil
}

// "/project", "/location", "/inventory" POST requests
// returns true if successful otherwise it logs the error and returns false
func SetModelFire(model model.Model, modelName string) bool {

	_, err := global_client.Collection(modelName).Doc(model.ID()).Set(global_ctx, model)
	if err != nil {
		c.ErrorLog.Println(err.Error())
		return false
	}

	return true
}

// wrapper for "/device" POST requests
func SetModelFireWrapper(model model.Model, modelName, foreign string) error {

	if model.FID() != "" {
		_, err := global_client.Collection(foreign).Doc(model.FID()).Get(global_ctx)
		if err != nil {
			c.ErrorLog.Println(err.Error())
			if status.Code(err) == codes.NotFound {
				return errors.New(foreign + " " + model.FID() + " was not found")
			}
			return err
		}
	}

	if SetModelFire(model, modelName) {
		return nil
	}

	return errors.New("database error")
}

func UpdatePingFire(ping *model.Ping) (*model.Ping, error) {

	devSnap, _ := global_client.Collection("device").Doc(ping.Id).Get(global_ctx)
	pingSnap, _ := devSnap.DataAt("lastPing")

	_, err := global_client.Collection("device").Doc(ping.Id).Update(global_ctx, []firestore.Update{
		{
			Path:  "lastPing",
			Value: ping,
		},
	})
	if err != nil {
		c.ErrorLog.Println(err.Error())
		return nil, errors.New("device to update not found")
	}

	if val, ok := pingSnap.(map[string]interface{}); ok {
		err := mapstructure.Decode(val, &ping)

		ping.Timestamp = val["timestamp"].(time.Time)

		if err != nil {
			c.ErrorLog.Println(err)
			c.WarningLog.Println("current ping returned instead of previous one")
			return ping, nil
		}
	}

	return ping, nil
}

func SetUserFire(user *model.UserCredentials) error {

	_, exists := global_client.Collection("user").Doc(user.Id).Get(global_ctx)

	if status.Code(exists) == codes.NotFound {
		_, err := global_client.Collection("user").Doc(user.Id).Set(global_ctx, *user)
		if err != nil {
			c.ErrorLog.Println(err.Error())
			return errors.New("database error")
		}
		return nil
	} else {
		c.ErrorLog.Printf("user %s already exists\n", user.Username)
		return errors.New("user already exists")
	}
}

package route

import (
	"encoding/json"
	"errors"
	"fmt"
	c "health/clog"
	"health/model"
	"health/network"
	"strings"
	"sync"

	"github.com/mitchellh/mapstructure"
)

func PostLocation(locationMap map[string]interface{}) (string, error) {
	var address = model.LocationAddress{}
	var location = model.Location{Address: address}

	if missing := validateLocation(locationMap); len(missing) != 0 {
		m := strings.Join(missing, ", ")
		c.ErrorLog.Printf("missing parameters: %s\n", m)
		return "", fmt.Errorf("failed to upload location: missing required parameter(s): %s", m)
	}

	if err := mapstructure.Decode(locationMap, &location); err != nil {
		c.ErrorLog.Println(err.Error())
		return "", errors.New("failed to upload location: decoding error")
	}

	if location.Id == "" {
		location.Id = generateUUID()
	}

	if err := network.SetModelFireWrapper(&location, "location", "project"); err != nil {
		return "", fmt.Errorf("failed to upload location %s", err.Error())
	}

	return location.Id, nil
}

func validateLocation(locationMap map[string]interface{}) []string {

	missing := validateParameters(model.LocationParameters, locationMap)
	for _, v := range missing {
		if v == "address" {
			return missing
		}
	}
	return append(missing, validateParameters(model.AddressParameters, locationMap["address"].(map[string]interface{}))...)
}

func DeleteLocation(locationId string, out *[]byte, wgProject *sync.WaitGroup) error {
	defer removeWaitGroup(wgProject)

	doc, err := network.GetSingleDocument("location", locationId)
	if err != nil {
		return err
	}

	jsonBytes, err := network.GetAllDocuments("device")
	if err != nil {
		return err
	}

	var device []model.Device
	_ = json.Unmarshal(jsonBytes, &device)

	wg := sync.WaitGroup{}
	for _, dev := range device {

		if dev.FID() == locationId {
			wg.Add(1)
			go DeleteDevice(dev.ID(), nil, &wg)
		}
	}
	wg.Wait()

	if out != nil {
		*out, _ = json.Marshal(doc)
	}
	return network.DeleteFire("location", locationId)
}

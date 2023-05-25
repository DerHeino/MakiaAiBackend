package route

import (
	json "encoding/json"
	"errors"
	"fmt"
	"health/model"
	"health/network"
	"log"
	"strings"

	mapstructure "github.com/mitchellh/mapstructure"
)

var device model.Device

func PostDevice(deviceMap map[string]interface{}) (string, error) {
	defer clearModel(&device)

	if err := validateDevice(deviceMap); err != nil {
		log.Printf("missing parameters: %s", err.Error())
		return "", fmt.Errorf("missing required parameter(s): %s", err)
	}

	if err := mapstructure.Decode(deviceMap, &device); err != nil {
		log.Println(err.Error())
		return "", errors.New("failed to convert device")
	}

	if val, ok := deviceMap["lastPing"].(map[string]interface{}); ok {
		device.LastPing.Timestamp = controlTime(val["timestamp"].(string))
	}

	if err := network.SetDeviceFire(&device); err == nil {
		jsonStr, err := json.Marshal(&device)
		if err != nil {
			log.Println(err.Error())
			return "", errors.New("failed to convert uploaded device")
		}
		return string(jsonStr), nil
	} else {
		log.Println(err.Error())
		return "", errors.New("failed to upload device\nreason: " + err.Error())
	}
}

func validateDevice(deviceMap map[string]interface{}) error {

	missing := CountParameters(model.DeviceParameters, deviceMap)

	if ping, ok := deviceMap["lastPing"].(map[string]interface{}); ok {

		if err := validatePing(ping); err != nil {
			missing = append(missing, err.Error())
			return errors.New(strings.Join(missing, ", "))
		}
	}

	for _, p := range missing {
		if p == "_id" {
			deviceMap["_id"] = generateUUID()
		}
		if p == "name" || p == "locationId" {
			return errors.New(strings.Join(missing, ", "))
		}
	}

	return nil
}

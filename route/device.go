package route

import (
	json "encoding/json"
	"errors"
	"fmt"
	c "health/clog"
	"health/model"
	"health/network"
	"strings"

	mapstructure "github.com/mitchellh/mapstructure"
)

var device model.Device

//var ping model.Ping

func PostDevice(deviceMap map[string]interface{}) (string, error) {
	defer clearModel(&device)

	if err := validateDevice(deviceMap); err != nil {
		c.ErrorLog.Printf("missing parameters: %s\n", err.Error())
		return "", fmt.Errorf("failed to upload device: missing required parameter(s): %s", err)
	}

	if err := mapstructure.Decode(deviceMap, &device); err != nil {
		c.ErrorLog.Println(err.Error())
		return "", errors.New("failed to upload device: decoding error")
	}

	if val, ok := deviceMap["lastPing"].(map[string]interface{}); ok {
		device.LastPing.Timestamp = controlTime(val["timestamp"].(string))
	}

	if err := network.SetDeviceFire(&device); err == nil {
		jsonStr, err := json.Marshal(&device)
		if err != nil {
			c.WarningLog.Println(err.Error())
			return "", errors.New("uploaded")
		}
		return string(jsonStr), nil
	} else {
		return "", fmt.Errorf("failed to upload device: %s", err.Error())
	}
}

func validateDevice(deviceMap map[string]interface{}) error {

	id, missing := CountParameters(model.DeviceParameters, deviceMap)

	if ping, ok := deviceMap["lastPing"].(map[string]interface{}); ok {

		if err := validatePing(ping); err != nil {
			missing = append(missing, err.Error())
			return errors.New(strings.Join(missing, ", "))
		}
	}

	if !id {
		deviceMap["_id"] = generateUUID()
	}

	if len(missing) > 0 {
		return errors.New(strings.Join(missing, ", "))
	}

	return nil
}

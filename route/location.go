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

var location model.Location

func PostLocation(locationMap map[string]interface{}) (string, error) {
	defer clearModel(&location)

	if err := validateLocation(locationMap); err != nil {
		c.ErrorLog.Printf("missing parameters: %s\n", err.Error())
		return "", fmt.Errorf("failed to upload location: missing required parameter(s): %s", err)
	}

	if err := mapstructure.Decode(locationMap, &location); err != nil {
		c.ErrorLog.Println(err.Error())
		return "", errors.New("failed to upload location: decoding error")
	}

	if network.SetLocationFire(&location) {
		jsonStr, err := json.Marshal(&location)
		if err != nil {
			c.ErrorLog.Println(err.Error())
			return "", errors.New("uploaded")
		}
		return string(jsonStr), nil
	} else {
		return "", errors.New("failed to upload location")
	}
}

func validateLocation(locationMap map[string]interface{}) error {

	id, missing := CountParameters(model.LocationParameters, locationMap)

	if add, ok := locationMap["address"].(map[string]interface{}); ok {
		_, missingAddress := CountParameters(model.AddressParameters, add)
		if len(missingAddress) > 0 {
			missing = append(missing, missingAddress...)
			return errors.New(strings.Join(missing, ", "))
		}
	} else {
		return errors.New(strings.Join(missing, ", "))
	}

	if !id {
		locationMap["_id"] = generateUUID()
	}

	if len(missing) > 0 {
		return errors.New(strings.Join(missing, ", "))
	}

	return nil
}

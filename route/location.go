package route

import (
	json "encoding/json"
	"errors"
	"health/model"
	"health/network"
	"log"
	"strings"

	mapstructure "github.com/mitchellh/mapstructure"
)

var location model.Location

func PostLocation(locationMap map[string]interface{}) (string, error) {
	defer clearModel(&location)

	if err := validateLocation(locationMap); err != nil {
		log.Println(err.Error())
		return "", err
	}

	if err := mapstructure.Decode(locationMap, &location); err != nil {
		log.Println(err.Error())
		return "", errors.New("failed to convert location")
	}

	if network.SetLocationFire(&location) {
		jsonStr, err := json.Marshal(&location)
		if err != nil {
			log.Println(err.Error())
			return "", errors.New("failed to convert uploaded location")
		}
		return string(jsonStr), nil
	} else {
		return "", errors.New("failed to upload location")
	}
}

func validateLocation(locationMap map[string]interface{}) error {

	missing := CountParameters(model.LocationParameters, locationMap)

	if add, ok := locationMap["address"].(map[string]interface{}); ok {
		missingAddress := CountParameters(model.AddressParameters, add)
		if len(missingAddress) > 0 {
			missing = append(missing, missingAddress...)
			return errors.New("missing parameter(s): " + strings.Join(missing, ", "))
		}
	} else {
		return errors.New("missing parameter(s): " + strings.Join(missing, ", "))
	}

	for _, p := range missing {
		if p == "_id" {
			locationMap["_id"] = generateUUID()
		}
		if p == "name" || p == "address" {
			return errors.New("missing required parameter(s): " + strings.Join(missing, ", "))
		}
	}

	return nil
}

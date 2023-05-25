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

var ping model.Ping

func PostPing(pingMap map[string]interface{}) (string, error) {
	defer clearModel(&ping)

	if err := validatePing(pingMap); err != nil {
		log.Printf("missing parameters: %s\n", err.Error())
		return "", fmt.Errorf("missing parameters: %s", err)
	}

	if err := mapstructure.Decode(pingMap, &ping); err != nil {
		log.Println(err.Error())
		return "", err
	}

	if val, ok := pingMap["timestamp"]; ok {
		ping.Timestamp = controlTime(val.(string))
	}

	if val, err := network.UpdatePingFire(&ping); err == nil {
		jsonStr, err := json.Marshal(&val)
		if err != nil {
			log.Println(err.Error())
			return "", errors.New("failed to convert uploaded device")
		}
		return string(jsonStr), nil
	} else {
		log.Println(err.Error())
		return "", errors.New("failed to update Ping\nreason: " + err.Error())
	}
}

func validatePing(pingMap map[string]interface{}) error {

	missing := CountParameters(model.PingParameters, pingMap)

	if len(missing) > 0 {
		return errors.New(strings.Join(missing, ", "))
	}

	if status, ok := pingMap["status"].(string); ok {
		if err := validateStatus(status); err != nil {
			return err
		}
	}

	return nil
}

func validateStatus(status string) error {

	for _, s := range model.DeviceStatus {
		if status == s {
			return nil
		}
	}
	return errors.New("invalid status")
}

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

var ping model.Ping

func PostPing(pingMap map[string]interface{}) (string, error) {
	defer clearModel(&ping)

	if err := validatePing(pingMap); err != nil {
		c.ErrorLog.Printf("missing parameters: %s\n", err.Error())
		return "", fmt.Errorf("failed to update ping missing parameters: %s", err)
	}

	if err := mapstructure.Decode(pingMap, &ping); err != nil {
		c.ErrorLog.Println(err.Error())
		return "", errors.New("failed to update ping: decoding error")
	}

	if val, ok := pingMap["timestamp"]; ok {
		ping.Timestamp = controlTime(val.(string))
	}

	if val, err := network.UpdatePingFire(&ping); err == nil {
		jsonStr, err := json.Marshal(&val)
		if err != nil {
			c.ErrorLog.Println(err.Error())
			return "", errors.New("return error")
		}
		return string(jsonStr), nil
	} else {
		c.ErrorLog.Println(err.Error())
		return "", errors.New("failed to update ping: " + err.Error())
	}
}

func validatePing(pingMap map[string]interface{}) error {

	_, missing := CountParameters(model.PingParameters, pingMap)

	if len(missing) > 0 {
		return errors.New(strings.Join(missing, ", "))
	}

	if status, ok := pingMap["status"].(string); ok {
		return validateStatus(status)
	}

	return nil
}

func validateStatus(status string) error {

	for _, s := range model.DeviceStatus {
		if status == s {
			return nil
		}
	}
	return errors.New("invalid DeviceStatus")
}

package route

import (
	"encoding/json"
	"errors"
	"fmt"
	c "health/clog"
	"health/model"
	"health/network"
	"strings"

	"github.com/mitchellh/mapstructure"
)

func PostPing(pingMap map[string]interface{}) (string, error) {
	var ping = model.Ping{}

	if missing := validateParameters(model.PingParameters, pingMap); len(missing) != 0 {
		m := strings.Join(missing, ", ")
		c.ErrorLog.Printf("missing parameters: %s\n", m)
		return "", fmt.Errorf("failed to update ping missing parameters: %s", m)
	}

	if err := validatePingStatus(pingMap["status"].(string)); err != nil {
		return "", err
	}

	if err := mapstructure.Decode(pingMap, &ping); err != nil {
		c.ErrorLog.Println(err.Error())
		return "", errors.New("failed to update ping: decoding error")
	}

	if val, ok := pingMap["timestamp"]; ok {
		ping.Timestamp = decodeTimeWrapper(val)
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

func validatePingStatus(status string) error {

	for _, s := range model.DeviceStatus {
		if status == s {
			return nil
		}
	}
	return errors.New("invalid DeviceStatus")
}

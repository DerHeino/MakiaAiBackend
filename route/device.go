package route

import (
	"errors"
	"fmt"
	bg "health/background"
	c "health/clog"
	"health/model"
	"health/network"
	"image"
	"strings"

	"github.com/mitchellh/mapstructure"
)

func PostDevice(deviceMap map[string]interface{}) (string, error) {
	var device = model.Device{}

	if missing := validateParameters(model.DeviceParameters, deviceMap); len(missing) != 0 {
		m := strings.Join(missing, ", ")
		c.ErrorLog.Printf("missing parameters: %s\n", m)
		return "", fmt.Errorf("failed to upload device: missing required parameter(s): %s", m)
	}

	if err := mapstructure.Decode(deviceMap, &device); err != nil {
		c.ErrorLog.Println(err.Error())
		return "", errors.New("failed to upload device: decoding error")
	}

	if device.Id == "" {
		device.Id = generateUUID()
	}

	if err := network.SetModelFireWrapper(&device, "device", "location"); err == nil {
		if _, ok := deviceMap["lastPing"].(map[string]interface{}); ok {
			_, err := PostPing(deviceMap["lastPing"].(map[string]interface{}))
			if err != nil {
				return device.Id, err
			}
			return device.Id, nil
		} else {
			return device.Id, nil
		}

	} else {
		return "", fmt.Errorf("failed to upload device: %s", err.Error())
	}
}

func PostImage(deviceId string, image *image.Image) bool {
	devMap := bg.GetDeviceMap()
	return devMap.Add(deviceId, image)
}

func GetImage(deviceId string) *image.Image {
	devMap := bg.GetDeviceMap()
	return devMap.Get(deviceId)
}

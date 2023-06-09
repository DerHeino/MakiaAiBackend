package route

import (
	"encoding/json"
	"errors"
	"fmt"
	bg "health/background"
	c "health/clog"
	"health/model"
	"health/network"
	"image"
	"strings"
	"sync"

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

	if device.Status == nil {
		device.Status = &model.DeviceStatus[2]
	}

	if err := network.SetModelFireWrapper(&device, "device", "location"); err == nil {
		devMap := bg.GetDeviceMap()

		if _, ok := deviceMap["lastPing"].(map[string]interface{}); ok {
			_, err := PostPing(deviceMap["lastPing"].(map[string]interface{}))
			if err != nil {
				devMap.AddDevice(device.Id)
				return device.Id, err
			}
			devMap.AddDevice(device.Id)
			return device.Id, nil
		} else {
			devMap.AddDevice(device.Id)
			return device.Id, nil
		}

	} else {
		return "", fmt.Errorf("failed to upload device: %s", err.Error())
	}
}

func DeleteDevice(deviceId string, out *[]byte, wgLocation *sync.WaitGroup) error {
	defer removeWaitGroup(wgLocation)

	doc, err := network.GetSingleDocument("device", deviceId)
	if err != nil {
		return err
	}

	jsonBytes, err := network.GetAllDocuments("inventory")
	if err != nil {
		return err
	}

	var inventory []model.Inventory
	_ = json.Unmarshal(jsonBytes, &inventory)

	wg := sync.WaitGroup{}
	for _, inv := range inventory {

		if inv.FID() == deviceId {
			wg.Add(1)
			go DeleteInventory(inv.ID(), nil, &wg)
		}
	}
	wg.Wait()

	bg.GetDeviceMap().Delete(deviceId)
	if out != nil {
		*out, _ = json.Marshal(doc)
	}
	return network.DeleteFire("device", deviceId)
}

func PostImage(deviceId string, image *image.Image) bool {
	return bg.GetDeviceMap().AddImage(deviceId, image)
}

func GetImage(deviceId string) *image.Image {
	return bg.GetDeviceMap().Get(deviceId)
}

package route

import (
	"fmt"
	"health/model"
	"health/network"
	"image"
)

var deviceList map[string]model.LocalDevice

func getDeviceList() {
	deviceList = network.DeviceList
}

func PostImage(deviceId string, image *image.Image) bool {
	getDeviceList()

	fmt.Println(deviceList, "ri")
	if entry, ok := deviceList[deviceId]; ok {
		entry.Image = image

		deviceList[deviceId] = entry
		return true
	}

	return false
}

func GetImage(deviceId string) *image.Image {
	getDeviceList()

	if _, ok := deviceList[deviceId]; ok {
		if deviceList[deviceId].Image != nil {
			return deviceList[deviceId].Image
		} else {
			return nil
		}
	}

	return nil
}

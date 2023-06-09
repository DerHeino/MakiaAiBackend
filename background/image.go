package background

import (
	c "health/clog"
	"health/network"
	"image"
	"sync"
	"time"
)

type localImage struct {
	image   *image.Image
	lastAck int64
}

type DeviceMap struct {
	devices map[string]*localImage
	l       sync.Mutex
}

var devMap *DeviceMap

func RunDevMap() {

	devMap = &DeviceMap{devices: make(map[string]*localImage)}

	go func() {
		var now time.Time
		ticker := time.NewTicker(time.Minute * 5)
		for ; true; now = <-ticker.C {
			if now.IsZero() {
				now = time.Now()
			}
			fireDevices := network.GetAllDeviceIDs()
			devMap.l.Lock()
			for _, firedev := range fireDevices {
				if entry, ok := devMap.devices[firedev]; ok {
					entry.lastAck = now.Unix()
				} else {
					devMap.devices[firedev] = &localImage{lastAck: now.Unix()}
				}
			}
			for id, entry := range devMap.devices {
				if now.Unix()-entry.lastAck > 60 {
					c.InfoLog.Printf("backround: devmap delete expired device: %s", id)
					devMap.deleteImage(id)
				}
			}
			devMap.l.Unlock()
		}
	}()
}

func GetDeviceMap() *DeviceMap {
	return devMap
}

func (devMap *DeviceMap) Len() int {
	return len(devMap.devices)
}

func (devMap *DeviceMap) AddDevice(id string) (b bool) {
	devMap.l.Lock()
	if _, ok := devMap.devices[id]; !ok {
		devMap.devices[id] = &localImage{lastAck: time.Now().Unix()}
		b = true
	}
	devMap.l.Unlock()
	return
}

func (devMap *DeviceMap) AddImage(id string, img *image.Image) (b bool) {
	devMap.l.Lock()
	if entry, ok := devMap.devices[id]; ok {
		entry.image = img
		entry.lastAck = time.Now().Unix()
		b = true
	} else {
		b = false
	}
	devMap.l.Unlock()
	return
}

func (devMap *DeviceMap) Get(id string) (img *image.Image) {
	devMap.l.Lock()
	if entry, ok := devMap.devices[id]; ok {
		img = entry.image
	}
	devMap.l.Unlock()
	return
}

func (devMap *DeviceMap) Delete(id string) {
	devMap.l.Lock()
	devMap.deleteImage(id)
	devMap.l.Unlock()
}

func (devMap *DeviceMap) deleteImage(id string) {
	if entry, ok := devMap.devices[id]; ok {
		entry.image = nil
	}
	delete(devMap.devices, id)
}

package route

import (
	"encoding/json"
	"errors"
	"fmt"
	c "health/clog"
	"health/model"
	"health/network"
	"strings"
	"sync"

	"github.com/mitchellh/mapstructure"
)

func PostInventory(inventoryMap map[string]interface{}) (string, error) {
	var inventory = model.Inventory{}

	if missing := validateParameters(model.InventoryParameters, inventoryMap); len(missing) != 0 {
		m := strings.Join(missing, ", ")
		c.ErrorLog.Printf("missing parameters: %s\n", m)
		return "", fmt.Errorf("failed to upload inventory: missing required parameter(s): %s", m)
	}

	if err := mapstructure.Decode(inventoryMap, &inventory); err != nil {
		c.ErrorLog.Println(err.Error())
		return "", errors.New("failed to upload inventory: decoding error")
	}

	if inventory.Id == "" {
		inventory.Id = generateUUID()
	}

	if val, ok := inventoryMap["buyDate"]; ok {
		inventory.BuyDate = decodeTime(val.(string))
	}

	if val, ok := inventoryMap["montageDate"]; ok {
		inventory.MontageDate = decodeTime(val.(string))
	}

	if err := network.SetModelFireWrapper(&inventory, "inventory", "device"); err != nil {
		return "", fmt.Errorf("failed to upload inventory: %s", err.Error())
	}

	return inventory.Id, nil
}

func DeleteInventory(id string, out *[]byte, wgDevice *sync.WaitGroup) error {
	defer removeWaitGroup(wgDevice)

	// firestore delete alone does not notify if document does not exists, so manually
	doc, err := network.GetSingleDocument("inventory", id)
	if err != nil {
		return err
	}
	if out != nil {
		*out, _ = json.Marshal(doc)
	}
	return network.DeleteFire("inventory", id)
}

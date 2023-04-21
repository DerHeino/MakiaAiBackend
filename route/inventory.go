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

var inventory model.Inventory

func PostInventory(inventoryMap map[string]interface{}) (string, error) {
	defer clearModel(&inventory)

	if err := validateInventory(inventoryMap); err != nil {
		c.ErrorLog.Printf("missing parameters: %s\n", err.Error())
		return "", fmt.Errorf("failed to upload inventory: missing required parameter(s): %s", err)
	}

	if err := mapstructure.Decode(inventoryMap, &inventory); err != nil {
		c.ErrorLog.Println(err.Error())
		return "", errors.New("failed to upload inventory: decoding error")
	}

	if val, ok := inventoryMap["buyDate"]; ok {
		inventory.BuyDate = controlTime(val.(string))
	}

	if val, ok := inventoryMap["montageDate"]; ok {
		inventory.MontageDate = controlTime(val.(string))
	}

	if network.SetInventoryFire(&inventory) {
		jsonStr, err := json.Marshal(&inventory)
		if err != nil {
			c.ErrorLog.Println(err.Error())
			return "", errors.New("uploaded")
		}
		return string(jsonStr), nil
	} else {
		return "", errors.New("failed to upload inventory")
	}
}

func validateInventory(inventoryMap map[string]interface{}) error {

	id, missing := CountParameters(model.InventoryParameters, inventoryMap)

	if !id {
		inventoryMap["_id"] = generateUUID()
	}

	if len(missing) > 0 {
		return errors.New(strings.Join(missing, ", "))
	}

	return nil
}

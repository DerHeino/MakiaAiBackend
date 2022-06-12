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

var inventory model.Inventory

func PostInventory(inventoryMap map[string]interface{}) (string, error) {
	defer clearModel(&inventory)

	if err := validateInventory(inventoryMap); err != nil {
		log.Println(err.Error())
		return "", err
	}

	if err := mapstructure.Decode(inventoryMap, &inventory); err != nil {
		log.Println(err.Error())
		return "", errors.New("failed to convert inventory")
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
			log.Println(err.Error())
			return "", errors.New("failed to convert uploaded inventory")
		}
		return string(jsonStr), nil
	} else {
		return "", errors.New("failed to upload inventory")
	}
}

func validateInventory(inventoryMap map[string]interface{}) error {

	missing := CountParameters(model.InventoryParameters, inventoryMap)

	for _, p := range missing {
		if p == "_id" {
			inventoryMap["_id"] = generateUUID()
		}
		if p == "name" {
			return errors.New("missing required parameter: " + strings.Join(missing, ", "))
		}
	}

	return nil
}

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

var project model.Project

func PostProject(projectMap map[string]interface{}) (string, error) {
	defer clearModel(&project)

	if err := validateProject(projectMap); err != nil {
		c.ErrorLog.Printf("missing parameters: %s\n", err.Error())
		return "", fmt.Errorf("failed to upload project: missing required parameter(s): %s", err)
	}

	if err := mapstructure.Decode(projectMap, &project); err != nil {
		c.ErrorLog.Println(err.Error())
		return "", errors.New("failed to upload project: decoding error")
	}

	if network.SetProjectFire(&project) {
		jsonStr, err := json.Marshal(&project)
		if err != nil {
			c.ErrorLog.Println(err.Error())
			return "", errors.New("uploaded")
		}
		return string(jsonStr), nil
	} else {
		return "", errors.New("failed to upload project")
	}
}

func validateProject(projectMap map[string]interface{}) error {

	id, missing := CountParameters(model.ProjectParameters, projectMap)

	if !id {
		projectMap["_id"] = generateUUID()
	}

	if len(missing) > 0 {
		return errors.New(strings.Join(missing, ", "))
	}

	return nil
}

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

var project model.Project

func PostProject(projectMap map[string]interface{}) (string, error) {
	defer clearModel(&project)

	if err := validateProject(projectMap); err != nil {
		return "", err
	}

	if err := mapstructure.Decode(projectMap, &project); err != nil {
		log.Println(err.Error())
		return "", errors.New("failed to convert project")
	}

	if network.SetProjectFire(&project) {
		jsonStr, err := json.Marshal(&project)
		if err != nil {
			log.Println(err.Error())
			return "", errors.New("failed to convert uploaded project")
		}
		return string(jsonStr), nil
	} else {
		return "", errors.New("failed to upload project")
	}
}

func validateProject(projectMap map[string]interface{}) error {

	missing := CountParameters(model.ProjectParameters, projectMap)

	if len(missing) > 0 {
		return errors.New("missing required parameter(s): " + strings.Join(missing, ", "))
	}

	return nil
}

package route

import (
	"errors"
	"fmt"
	c "health/clog"
	"health/model"
	"health/network"
	"strings"

	"github.com/mitchellh/mapstructure"
)

func PostProject(projectMap map[string]interface{}) (string, error) {
	var project = model.Project{}

	if missing := validateParameters(model.ProjectParameters, projectMap); len(missing) != 0 {
		m := strings.Join(missing, ", ")
		c.ErrorLog.Printf("missing parameters: %s\n", m)
		return "", fmt.Errorf("failed to upload project: missing required parameter(s): %s", m)
	}

	if err := mapstructure.Decode(projectMap, &project); err != nil {
		c.ErrorLog.Println(err.Error())
		return "", errors.New("failed to upload project: decoding error")
	}

	if project.Id == "" {
		project.Id = generateUUID()
	}

	if network.SetModelFire(&project, "project") {
		return project.Id, nil
	} else {
		return "", errors.New("failed to upload project")
	}
}

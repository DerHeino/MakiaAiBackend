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

func DeleteProject(projectId string, out *[]byte) error {

	doc, err := network.GetSingleDocument("project", projectId)
	if err != nil {
		return err
	}

	jsonBytes, err := network.GetAllDocuments("location")
	if err != nil {
		return err
	}

	var location []model.Location
	_ = json.Unmarshal(jsonBytes, &location)

	wg := sync.WaitGroup{}
	for _, loc := range location {

		if loc.FID() == projectId {
			wg.Add(1)
			go DeleteLocation(loc.ID(), nil, &wg)
		}
	}
	wg.Wait()

	*out, _ = json.Marshal(doc)
	return network.DeleteFire("project", projectId)
}

package route

import (
	"fmt"
	c "health/clog"
	"reflect"
	"time"

	"github.com/google/uuid"
)

// Generates and returns UUID
func generateUUID() string {
	return uuid.New().String()
}

// Checks map for missing parameters
//
// Certain parameters such as "_id" can be omited in the request (in that case the backend will generate one)
// for required parameters an error should be returned to the client
//
// returns an array containing which parameters are missing from the map
func CountParameters(parameter []string, model map[string]interface{}) (bool, []string) {
	var id = true
	var missing = []string{}

	for _, p := range parameter {

		if val, ok := model[p]; !ok {
			if p == "_id" {
				id = false
			} else {
				missing = append(missing, p)
			}
		} else if val == reflect.Zero(reflect.TypeOf(val)) {
			if p == "_id" {
				id = false
			} else {
				missing = append(missing, p)
			}
		}
	}

	return id, missing
}

// Returns a pointer containing the time read from value in RFC3339 standard
func controlTime(value string) time.Time {
	var err error

	updateTime, err := time.Parse(time.RFC3339, value)
	if err != nil {
		c.WarningLog.Println(err.Error())
		return updateTime
	}

	return updateTime
}

// Clears entire model struct of previous values.
// Parameter has to be a pointer otherwise it will panic #
func clearModel(model interface{}) {
	if reflect.ValueOf(model).Kind() == reflect.Struct {
		value := reflect.ValueOf(model).Elem()
		value.Set(reflect.Zero(value.Type()))
		fmt.Printf("clearModel %p\n\t%+v\n", model, model)
	}
}

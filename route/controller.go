package route

import (
	"log"
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
func CountParameters(parameter []string, model map[string]interface{}) []string {
	var missing = []string{}

	for _, p := range parameter {

		if val, ok := model[p]; !ok {
			missing = append(missing, p)
		} else if val == reflect.Zero(reflect.TypeOf(val)) {
			missing = append(missing, p)
		}
	}

	return missing
}

// Returns a pointer containing the time read from value in RFC3339 standard
func controlTime(value string) time.Time {
	var err error

	updateTime, err := time.Parse(time.RFC3339, value)
	if err != nil {
		log.Println(err.Error(), "controlTime")
		return updateTime
	}

	return updateTime
}

// Clears entire model struct of previous values.
// Parameter has to be a pointer otherwise it will panic #
func clearModel(model interface{}) {
	value := reflect.ValueOf(model).Elem()
	value.Set(reflect.Zero(value.Type()))
}

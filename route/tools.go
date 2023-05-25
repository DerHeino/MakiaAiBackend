package route

import (
	c "health/clog"
	"reflect"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// Generates and returns UUID
func generateUUID() string {
	return uuid.New().String()
}

// Checks map for missing parameters
// returns an array containing which parameters are missing from the map
func validateParameters(parameter []string, model map[string]interface{}) []string {
	var missing = []string{}

	for _, p := range parameter {

		if val, ok := model[p]; !ok || (val == reflect.Zero(reflect.TypeOf(val))) {
			missing = append(missing, p)
		}
	}

	return missing
}

// Converts timestamp (UNIX and RFC3339) string to time struct
func decodeTime(input string) time.Time {

	if utime, err := strconv.ParseInt(input, 10, 64); err == nil {
		return time.Unix(utime, 0)
	} else if rfctime, err := time.Parse(time.RFC3339, input); err == nil {
		return rfctime
	} else {
		c.WarningLog.Println(err.Error())
		return time.Time{}
	}
}

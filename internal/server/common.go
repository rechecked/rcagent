
package server

import (
	"encoding/json"
)

func ConvertToJson(i interface{}, pretty bool) ([]byte, error) {
    if pretty {
        return json.MarshalIndent(i, "", "    ")
    }
    return json.Marshal(i)
}

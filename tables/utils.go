package config

import "fmt"

func IsArray(value interface{}) bool {
	_, isArray := value.([]interface{})
	return isArray
}

func IsMap(value interface{}) bool {
	_, isMap := value.(map[string]interface{})
	return isMap
}

func ValueToType(value interface{}) *string {
	dataType := fmt.Sprintf("%T", value)
	if dataType == "<nil>" {
		return nil
	}
	return &dataType
}

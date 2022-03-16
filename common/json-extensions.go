package common

import "encoding/json"

func ReadJson[T any](bytes []byte) (*T, error) {
	var value T
	if err := json.Unmarshal(bytes, &value); err != nil {
		return nil, err
	}
	return &value, nil
}

func ToJson(data interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

package utils

import (
	"encoding/json"
	"io"
)

func GetRequestBodyJSON(requestBody io.ReadCloser, varToDecodeTo interface{}) error {
	jsonDecoder := json.NewDecoder(requestBody)
	return jsonDecoder.Decode(varToDecodeTo)
}

func EncodeJSONResponse(dataToEncode interface{}) (string, error) {
	jsonBytes, err := json.Marshal(dataToEncode)

	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

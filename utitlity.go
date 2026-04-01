package main

import (
	"encoding/json"
	"io"
)

func jsonDecode[D any](data io.ReadCloser) (D, error) {
	decoder := json.NewDecoder(data)
	var formattedData D
	err := decoder.Decode(&formattedData)

	if err != nil {
		return formattedData, err
	}
	return formattedData, nil
}

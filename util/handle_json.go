package util

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
)

func JsonMarshalData(jsonData interface{}) ([]byte, error) {
	byteData, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}

	return byteData, nil
}

func JsonMarshalDataIndent(jsonData interface{}) ([]byte, error) {
	byteData, err := json.MarshalIndent(jsonData, "", "    ")
	if err != nil {
		return nil, err
	}

	return byteData, nil
}

func JsonUnmarshalData(jsonStruct interface{}, byteValue []byte) interface{} {
	json.Unmarshal(byteValue, &jsonStruct)

	return jsonStruct
}

func JsonUnmarshal(jsonStruct interface{}, jsonFilePath string) interface{} {
	jsonData, err := os.Open(jsonFilePath)
	if err != nil {
		LogErr(err)
	}
	byteValue, _ := io.ReadAll(jsonData)
	jsonStruct = JsonUnmarshalData(jsonStruct, byteValue)

	return jsonStruct
}

func SaveJsonPretty(jsonByte []byte, saveTxPath string) error {
	var prettyJson bytes.Buffer
	err := json.Indent(&prettyJson, jsonByte, "", "    ")
	if err != nil {
		LogErr(err)
	}

	err = os.WriteFile(saveTxPath, prettyJson.Bytes(), 0660)
	if err != nil {
		LogErr(err)
	}

	return nil
}

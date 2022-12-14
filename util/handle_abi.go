package util

import "os"

func AbiParsing(jsonFilePath string) (string, error) {
	f, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return "", err
	}
	return string(f), nil
}

type bytecodeParsingStruct struct {
}

func BytecodeParsing(jsonFilePath string) string {
	var bytecodeStruct bytecodeParsingStruct
	jsonData := JsonUnmarshal(bytecodeStruct, jsonFilePath)
	bytecode := jsonData.(map[string]interface{})["object"].(string)

	return bytecode
}

// For invoke(as execute) contract, parameters are packed by using ABI.
func GetAbiPack(callName string, args ...interface{}) ([]byte, error) {
	contractAbi, err := XplaSolContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	var abiByteData []byte

	if args == nil {
		abiByteData, err = contractAbi.Pack(callName)
		if err != nil {
			return nil, err
		}
	} else {
		abiByteData, err = contractAbi.Pack(callName, args...)
		if err != nil {
			return nil, err
		}
	}

	return abiByteData, nil
}

// After call(as query) solidity contract, the response of chain is unpacked by ABI.
func GetAbiUnpack(callName string, data []byte) ([]interface{}, error) {
	contractAbi, err := XplaSolContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	unpacked, err := contractAbi.Unpack(callName, data)
	if err != nil {
		return nil, err
	}

	return unpacked, nil
}

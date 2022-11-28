package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/Moonyongjung/xpla.go/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

func GetAddrByPrivKey(priv cryptotypes.PrivKey) sdk.AccAddress {
	addr, err := sdk.AccAddressFromHex(priv.PubKey().Address().String())
	if err != nil {
		LogErr(err)
	}
	return addr
}

func GasLimitAdjustment(gasUsed uint64, gasAdjustment string) (string, error) {
	gasAdj, err := strconv.ParseFloat(gasAdjustment, 64)
	if err != nil {
		return "", err
	}
	return FromIntToString(int(gasAdj * float64(gasUsed))), nil
}

func GrpcUrlParsing(normalUrl string) string {
	if strings.Contains(normalUrl, "http://") || strings.Contains(normalUrl, "https://") {
		parsedUrl := strings.Split(normalUrl, "://")
		return parsedUrl[1]
	} else {
		return normalUrl
	}
}

func AbiParsing(jsonFilePath string) (string, error) {
	f, err := ioutil.ReadFile(jsonFilePath)
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

func SaveJsonPretty(jsonByte []byte, saveTxPath string) error {
	var prettyJson bytes.Buffer
	err := json.Indent(&prettyJson, jsonByte, "", "    ")
	if err != nil {
		LogErr(err)
	}

	err = ioutil.WriteFile(saveTxPath, prettyJson.Bytes(), 0660)
	if err != nil {
		LogErr(err)
	}

	return nil
}

func MulUint64(val1 uint64, val2 uint64) uint64 {
	return val1 * val2
}

func MulBigInt(val1 *big.Int, val2 *big.Int) *big.Int {
	result := big.NewInt(0)
	return result.Mul(val1, val2)
}

func FromBigIntToString(v *big.Int) string {
	return v.String()
}

func FromStringToBigInt(v string) (*big.Int, error) {
	n := big.NewInt(0)
	n, ok := n.SetString(v, 10)
	if !ok {
		return nil, LogErr("convert string to big int err")
	}
	return n, nil
}

func FromStringToUint64(value string) uint64 {
	number, _ := strconv.ParseUint(value, 10, 64)
	return number
}

func FromUint64ToString(value uint64) string {
	return strconv.Itoa(int(value))
}

func FromStringToInt64(value string) int64 {
	number, _ := strconv.ParseInt(value, 10, 64)
	return number
}

func FromStringToInt(value string) int {
	number, _ := strconv.Atoi(value)
	return number
}

func FromIntToString(value int) string {
	return strconv.Itoa(value)
}

func FromStringToByte20Address(address string) common.Address {
	var byte20Address common.Address
	if address[:2] == "0x" {
		address = address[2:]
	}
	byte20Address = common.HexToAddress(address)

	return byte20Address
}

func FromByte20AddressToCosmosAddr(address common.Address) (sdk.AccAddress, error) {
	var addrStr string
	if address.Hex()[:2] == "0x" {
		addrStr = address.Hex()[2:]
	} else {
		addrStr = address.Hex()
	}

	addr, err := sdk.AccAddressFromHex(addrStr)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

func FromStringHexToHash(hashString string) common.Hash {
	return common.HexToHash(hashString)
}

func ToString(value interface{}, defaultValue string) string {
	s := fmt.Sprintf("%v", value)
	s = s[1 : len(s)-1]
	str := strings.TrimSpace(s)
	if str == "" {
		return defaultValue
	} else {
		return str
	}
}

func ToTypeHexString(value string) string {
	if !strings.Contains(value, "0x") {
		return "0x" + value
	} else {
		return value
	}
}

func JsonMarshalData(jsonData interface{}) ([]byte, error) {
	byteData, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}

	return byteData, err
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
	byteValue, _ := ioutil.ReadAll(jsonData)
	jsonStruct = JsonUnmarshalData(jsonStruct, byteValue)

	return jsonStruct
}

func DenomAdd(amount string) string {
	if strings.Contains(amount, types.XplaDenom) {
		return amount
	} else {
		return amount + types.XplaDenom
	}
}

func DenomRemove(amount string) string {
	if strings.Contains(amount, types.XplaDenom) {
		returnAmount := strings.Split(amount, types.XplaDenom)
		return returnAmount[0]
	} else {
		return amount
	}
}

func ConvertEvmChainId(chainId string) (*big.Int, error) {
	conv1 := strings.Split(chainId, "_")
	conv2 := strings.Split(conv1[1], "-")
	returnChainId, err := FromStringToBigInt(conv2[0])
	if err != nil {
		return nil, err
	}
	return returnChainId, nil
}

func Bech32toValidatorAddress(validators []string) ([]sdk.ValAddress, error) {
	vals := make([]sdk.ValAddress, len(validators))
	for i, validator := range validators {
		addr, err := sdk.ValAddressFromBech32(validator)
		if err != nil {
			return nil, err
		}
		vals[i] = addr
	}
	return vals, nil
}

func MakeQueryLcdUrl(metadata string) string {
	return "/" + strings.Replace(metadata, "query.proto", "", -1)
}

func MakeQueryLabels(labels ...string) string {
	return strings.Join(labels, "/")
}

func LogInfo(log ...interface{}) {
	fmt.Println(ToString(log, ""))
}

func LogErr(log ...interface{}) error {
	return errors.New(ToString(log, ""))
}

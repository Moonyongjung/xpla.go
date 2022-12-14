package util

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/Moonyongjung/xpla.go/types/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

func FromBigIntToString(v *big.Int) string {
	return v.String()
}

func FromStringToBigInt(v string) (*big.Int, error) {
	n := big.NewInt(0)
	n, ok := n.SetString(v, 10)
	if !ok {
		return nil, LogErr(errors.ErrInvalidRequest, "convert string to big int err")
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

func FromInt64ToString(value int64) string {
	return strconv.FormatInt(value, 10)
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

func From0xHexStringToIBignt(hexString string) *big.Int {
	n := new(big.Int)
	n.SetString(hexString[2:], 16)

	return n
}

func FromStringToTypeHexString(value string) string {
	if !strings.Contains(value, "0x") {
		return "0x" + value
	} else {
		return value
	}
}

func ToString(value interface{}, defaultValue string) string {
	s := fmt.Sprintf("%v", value)
	if s == "" {
		return defaultValue
	} else {
		return s
	}
}

func ToStringTrim(value interface{}, defaultValue string) string {
	s := fmt.Sprintf("%v", value)
	s = s[1 : len(s)-1]
	str := strings.TrimSpace(s)
	if str == "" {
		return defaultValue
	} else {
		return str
	}
}

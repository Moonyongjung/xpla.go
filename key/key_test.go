package key

import (
	"testing"

	"github.com/Moonyongjung/xpla.go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const mnemonicMykey = "quote wasp mixture bench upper flame salmon century viable dilemma can squirrel inmate away moon trigger echo sure measure doll peace abstract language nature"

func TestNewMnemonic(t *testing.T) {
	_, err := NewMnemonic()
	assert.NoError(t, err)
}

func TestNewPrivKey(t *testing.T) {
	mnemonic, err := NewMnemonic()
	assert.NoError(t, err)

	// Only Secp256k1 is supported
	_, err = NewPrivKey(mnemonic)
	assert.NoError(t, err)
}

func TestConvertPrivKeyToBech32Addr(t *testing.T) {
	util.SetChainConfig()
	addrMykey := "xpla1ww42aq5v6886w0aggpkrh9pcqudr7qggje32qf"

	PrivateKey, err := NewPrivKey(mnemonicMykey)
	assert.NoError(t, err)

	addr, err := Bech32AddrString(PrivateKey)
	assert.NoError(t, err)
	require.Equal(t, addrMykey, addr)
}

func TestConvertPrivKeyToHexAddr(t *testing.T) {
	util.SetChainConfig()
	addrMykey := "73AAAE828CD1CFA73FA8406C3B9438071A3F0108"

	PrivateKey, err := NewPrivKey(mnemonicMykey)
	assert.NoError(t, err)

	addr := HexAddrString(PrivateKey)
	require.Equal(t, addrMykey, addr)
}

func TestEncryptDecryptPrivKeyArmor(t *testing.T) {
	util.SetChainConfig()
	PrivateKey, err := NewPrivKey(mnemonicMykey)
	assert.NoError(t, err)

	armor1 := EncryptArmorPrivKey(PrivateKey, DefaultEncryptPassphrase)
	armor2 := EncryptArmorPrivKeyWithoutPassphrase(PrivateKey)

	pk1, algo1, err := UnarmorDecryptPrivKey(armor1, DefaultEncryptPassphrase)
	assert.NoError(t, err)
	pk2, algo2, err := UnarmorDecryptPrivKeyWithoutPassphrase(armor2)
	assert.NoError(t, err)

	require.Equal(t, pk1, pk2)
	require.Equal(t, algo1, algo2)
}

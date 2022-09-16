package key

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewMnemonic(t *testing.T) {
	_, err := NewMnemonic()
	assert.NoError(t, err)
}

func Test_NewPrivKey(t *testing.T) {
	mnemonic, err := NewMnemonic()
	assert.NoError(t, err)

	// Only Secp256k1 is supported
	_, err = NewPrivKey(mnemonic)
	assert.NoError(t, err)
}

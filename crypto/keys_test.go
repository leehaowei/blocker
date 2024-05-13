package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivateKey(t *testing.T) {
	privKey := GeneratePrivateKey()
	assert.Equal(t, len(privKey.Bytes()), privKeyLen)

	pubKey := privKey.Public()
	assert.Equal(t, len(pubKey.Bytes()), pubKeyLen)
}

func TestPrivateKeySign(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.Public()
	msg := []byte("foo bar baz")

	sig := privKey.Sing(msg)
	assert.True(t, sig.Verify(pubKey, msg))

	// Test with invalid msg
	assert.False(t, sig.Verify(pubKey, []byte("foo")))

	// test with invalid pubKey
	invalidPrivKey := GeneratePrivateKey()
	invalidPubkey := invalidPrivKey.Public()
	assert.False(t, sig.Verify(invalidPubkey, msg))
}

func TestNewPrivateKeyFromString(t *testing.T) {
	var (
		seed       = "d62f412da76a680b62007c7ab29739f22811b1ad2c22f9889e782e061327b4f0"
		privKey    = NewPrivateKeyFromString(seed)
		addressStr = "9028090886bee0ebf0894d94fb93e26594753952"
	)

	assert.Equal(t, privKeyLen, len(privKey.Bytes()))
	address := privKey.Public().Address()
	// fmt.Println(address)
	assert.Equal(t, addressStr, address.String())
}

func TestPublicKeyToAddress(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.Public()
	address := pubKey.Address()
	assert.Equal(t, addressLen, len(address.Byte()))
	// fmt.Println(address)
}

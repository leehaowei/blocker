package types

import (
	"testing"

	"github.com/leehaowei/blocker/util"
	"github.com/stretchr/testify/assert"

	"github.com/leehaowei/blocker/crypto"
)

func testSignBlock(t *testing.T) {
	var (
		block   = util.RandomBlock()
		prviKey = crypto.GeneratePrivateKey()
		pubKey  = prviKey.Public()
	)
	sig := SignBlock(prviKey, block)
	assert.Equal(t, 64, len(sig.Bytes()))
	assert.True(t, sig.Verify(pubKey, HashBlock(block)))

}

func TestHashBlock(t *testing.T) {
	block := util.RandomBlock()
	hash := HashBlock(block)
	// fmt.Println(hex.EncodeToString(hash))
	assert.Equal(t, 32, len(hash))
}

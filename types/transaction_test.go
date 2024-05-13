package types

import (
	"testing"

	"github.com/leehaowei/blocker/crypto"
	"github.com/leehaowei/blocker/proto"
	"github.com/leehaowei/blocker/util"
	"github.com/stretchr/testify/assert"
)

// Test case
// with balcance of 100 coins, send 5 coins to "AAA"
// 2 outputs
// 5 to the others
// 95 back to our address
func TestNewTransaction(t *testing.T) {
	fromPrivKey := crypto.GeneratePrivateKey()
	fromAddress := fromPrivKey.Public().Address().Byte()

	toPrivKey := crypto.GeneratePrivateKey()
	toAddress := toPrivKey.Public().Address().Byte()

	input := &proto.TxInput{
		PrevTxHash:   util.RandomHash(),
		PrevOutIndex: 0,
		PublicKey:    fromPrivKey.Public().Bytes(),
	}

	output1 := &proto.TxOutput{
		Amount:  5,
		Address: toAddress,
	}
	output2 := &proto.TxOutput{
		Amount:  95,
		Address: fromAddress,
	}

	tx := &proto.Transaction{
		Version: 1,
		Input:   []*proto.TxInput{input},
		Outputs: []*proto.TxOutput{output1, output2},
	}

	sig := SignTransaction(fromPrivKey, tx)
	input.Signature = sig.Bytes()

	assert.True(t, VerifyTransaction(tx))

	// fmt.Printf("%+v\n", tx)
}

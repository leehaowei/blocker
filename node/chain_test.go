package node

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/leehaowei/blocker/crypto"
	"github.com/leehaowei/blocker/proto"
	"github.com/leehaowei/blocker/types"
	"github.com/leehaowei/blocker/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func randomBlock(t *testing.T, chain *Chain) *proto.Block {
	privKey := crypto.GeneratePrivateKey()
	b := util.RandomBlock()
	prevBlock, err := chain.GetBlockByHeight(chain.Height())
	require.Nil(t, err)
	b.Header.PrevHash = types.HashBlock(prevBlock)
	types.SignBlock(privKey, b)
	return b
}

func TestNewChain(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTxStore())
	assert.Equal(t, 0, chain.Height())
	_, err := chain.GetBlockByHeight(0)
	assert.Nil(t, err)
}

func TestChainHeight(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTxStore())
	for i := 0; i < 100; i++ {
		b := randomBlock(t, chain)
		require.Nil(t, chain.AddBlock(b))
		require.Equal(t, chain.Height(), i+1)
	}
}

func TestAddBlock(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTxStore())

	for i := 0; i < 100; i++ {
		block := randomBlock(t, chain)
		blockHash := types.HashBlock(block)

		require.Nil(t, chain.AddBlock(block))

		fetchedBlock, err := chain.GetBlockByHash(blockHash)
		require.Nil(t, err)
		require.Equal(t, block, fetchedBlock)

		fetchedBlockBeyHeight, err := chain.GetBlockByHeight(i + 1)
		require.Nil(t, err)
		require.Equal(t, block, fetchedBlockBeyHeight)
	}
}

func TestAddBlockWithTx(t *testing.T) {
	var (
		chain     = NewChain(NewMemoryBlockStore(), NewMemoryTxStore())
		block     = randomBlock(t, chain)
		privKey   = crypto.NewPrivateKeyFromSeedStr(godSeed)
		recipient = crypto.GeneratePrivateKey().Public().Address().Byte()
	)

	ftt, err := chain.txStore.Get("5992c73a9240954b147a588a8e0290e1cb0f1743eb5a22e12936d192305c6e10")
	assert.Nil(t, err)

	inputs := []*proto.TxInput{
		{
			PrevTxHash:   types.HashTransaction(ftt),
			PrevOutIndex: 0,
			PublicKey:    privKey.Public().Bytes(),
		},
	}
	outputs := []*proto.TxOutput{
		{
			Amount:  100,
			Address: recipient,
		},
		{
			Amount:  900,
			Address: privKey.Public().Address().Byte(),
		},
	}
	tx := &proto.Transaction{
		Version: 1,
		Input:   inputs,
		Outputs: outputs,
	}

	sig := types.SignTransaction(privKey, tx)
	tx.Input[0].Signature = sig.Bytes()

	block.Transactions = append(block.Transactions, tx)
	require.Nil(t, chain.AddBlock(block))
	txHash := hex.EncodeToString(types.HashTransaction(tx))

	fetchedTx, err := chain.txStore.Get(txHash)
	assert.Nil(t, err)
	assert.Equal(t, tx, fetchedTx)

	// TODO: check if their is an UTXO that is unspent
	address := crypto.AddressFromByte(tx.Outputs[1].Address)
	key := fmt.Sprintf("%s_%s", address, txHash)

	utxo, err := chain.utxoStore.Get(key)
	assert.Nil(t, err)
	fmt.Println(utxo)
}

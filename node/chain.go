package node

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/leehaowei/blocker/crypto"
	"github.com/leehaowei/blocker/proto"
	"github.com/leehaowei/blocker/types"
)

const godSeed = "d233cc016102c74f8f449254d16264bda1b79c60cc24dd794842e4f8d36c9e36"

type HeaderList struct {
	headers []*proto.Header
}

func NewHeaderList() *HeaderList {
	return &HeaderList{
		headers: []*proto.Header{},
	}
}

func (list *HeaderList) Add(h *proto.Header) {
	list.headers = append(list.headers, h)
}

func (list *HeaderList) Get(index int) *proto.Header {
	if index > list.Height() {
		panic("index too high!")
	}
	return list.headers[index]
}

func (list *HeaderList) Height() int {
	return list.Len() - 1
}

// [A, B, C, D, E] Len = 5, height = 4
func (list *HeaderList) Len() int {
	return len(list.headers)
}

type UTXO struct {
	Hash     string
	OutIndex int
	Amount   int64
	Spent    bool
}

type Chain struct {
	txStore    TXStorer
	blockStore BlockStorer
	utxoStore  UTXOStorer
	headers    *HeaderList
}

func NewChain(bs BlockStorer, txStore TXStorer) *Chain {
	chain := &Chain{
		blockStore: bs,
		txStore:    txStore,
		utxoStore:  NewMemoryUTXOStore(),
		headers:    NewHeaderList(),
	}
	chain.addBlock(createGenisisBlock())
	return chain
}

func (c *Chain) Height() int {
	return c.headers.Height()
}

// Validate before actually adding the block
func (c *Chain) AddBlock(b *proto.Block) error {
	if err := c.ValidateBlock(b); err != nil {
		return err
	}
	return c.addBlock(b)
}

func (c *Chain) addBlock(b *proto.Block) error {
	c.headers.Add(b.Header)

	for _, tx := range b.Transactions {
		if err := c.txStore.Put(tx); err != nil {
			return err
		}
		hash := hex.EncodeToString(types.HashTransaction(tx))

		// address_txhash
		for i, output := range tx.Outputs {
			utxo := &UTXO{
				Hash:     hash,
				Amount:   output.Amount,
				OutIndex: i,
				Spent:    false,
			}
			address := crypto.AddressFromByte(output.Address)
			key := fmt.Sprintf("%s_%s", address, hash)
			if err := c.utxoStore.Put(key, utxo); err != nil {
				return err
			}
		}
	}
	// validation
	return c.blockStore.Put(b)
}

func (c *Chain) GetBlockByHash(hash []byte) (*proto.Block, error) {
	hashHex := hex.EncodeToString(hash)
	return c.blockStore.Get(hashHex)
}

func (c *Chain) GetBlockByHeight(height int) (*proto.Block, error) {
	if c.Height() < height {
		return nil, fmt.Errorf("given height (%d) too high - height (%d)", height, c.Height())
	}
	header := c.headers.Get(height)
	hash := types.HashHeader(header)
	return c.GetBlockByHash(hash)
}

func (c *Chain) ValidateBlock(b *proto.Block) error {
	// Validate the signature of the block
	if !types.VerifyBlock(b) {
		return fmt.Errorf("invalid block signature")
	}

	// Validate if the prevHash is the actually hash of the current block.
	currentBlock, err := c.GetBlockByHeight(c.Height())
	if err != nil {
		return err
	}
	hash := types.HashBlock(currentBlock)
	if !bytes.Equal(hash, b.Header.PrevHash) {
		return fmt.Errorf("invalid previous block hash")
	}

	for _, tx := range b.Transactions {
		if !types.VerifyTransaction(tx) {
			return fmt.Errorf("invalid tx signaure")
		}

		// TODO:
		// for _, input : range tx.Inputs {
		// }
	}

	return nil
}

func createGenisisBlock() *proto.Block {
	privKey := crypto.NewPrivateKeyFromSeedStr(godSeed)

	block := &proto.Block{
		Header: &proto.Header{
			Version: 1,
		},
	}
	types.SignBlock(privKey, block)

	tx := &proto.Transaction{
		Version: 1,
		Input:   []*proto.TxInput{},
		Outputs: []*proto.TxOutput{
			{
				Amount:  1000,
				Address: privKey.Public().Address().Byte(),
			},
		},
	}

	block.Transactions = append(block.Transactions, tx)
	types.SignBlock(privKey, block)

	return block
}

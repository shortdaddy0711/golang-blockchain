package main

import (
	"fmt"
	"bytes"
	"crypto/sha256"
)

// BlockChain structure
type BlockChain struct {
	blocks []*Block
}

// Block structure
type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
}

// DeriveHash method for Block structure
func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

// CreateBlock function to generate new block
func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash}
	block.DeriveHash()
	return block
}

// AddBlock method for BlockChain structure
func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.blocks[len(chain.blocks) - 1]
	new := CreateBlock(data, prevBlock.Hash)
	chain.blocks = append(chain.blocks, new)
}

// Genesis function to generate the Genesis block
func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

// InitBlockChain function to init blockchain with Genesis block
func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}

func main() {
	chain := InitBlockChain()

	chain.AddBlock("First Block after Genesis")
	chain.AddBlock("Second Block after Genesis")
	chain.AddBlock("Third Block after Genesis")

	for _, block := range chain.blocks {

	}
}
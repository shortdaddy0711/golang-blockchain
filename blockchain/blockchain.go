package blockchain


// BlockChain structure
type BlockChain struct {
	Blocks []*Block
}

// InitBlockChain function to init blockchain with Genesis block
func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}

// AddBlock method for BlockChain structure
func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks) - 1]
	new := CreateBlock(data, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, new)
}

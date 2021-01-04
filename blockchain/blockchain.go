package blockchain

import (
	"encoding/hex"
	"fmt"
	// "internal/unsafeheader"
	"os"
	"runtime"

	badger "github.com/dgraph-io/badger/v2"
)

const (
	dbPath = "./tmp/blocks"
	dbFile = "./tmp/blocks/MANIFEST"
	genesisData = "First Transaction from Genesis"
)

// BlockChain structure
type BlockChain struct {
	LastHash []byte
	Database *badger.DB
	// Blocks []*Block
}

// BlockChainIterator structure
type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

// DBexists function to check db exists or not
func DBexists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

// ContinueBlockChain function to add new block to existing blockchain
func ContinueBlockChain(address string) *BlockChain {
	if DBexists() == false {
		fmt.Println("No existing blockchain found, create one!")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh")) // retrieve last block of the blockchain
		Handle(err)
		lastHash, err = item.ValueCopy(nil)

		return err
	})
	Handle(err)

	chain := BlockChain{lastHash, db}
	return &chain
}

// InitBlockChain function to init blockchain with Genesis block
func InitBlockChain(address string) *BlockChain {
	var lastHash []byte

	if DBexists() {
		fmt.Println("Blockchain already exists")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(dbPath)

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		cbtx := CoinbaseTx(address, genesisData)
		genesis := Genesis(cbtx)
		fmt.Println("Genesis created")
		err = txn.Set(genesis.Hash, genesis.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), genesis.Hash) // save last hash to db

		lastHash = genesis.Hash // save last hash to memory

		return err
		// }
		// item, err := txn.Get([]byte("lh"))
		// Handle(err)
		// err = item.Value(func(val []byte) error {
		// 	lastHash = val
		// 	return nil
		// })
		// return err
	})

	Handle(err)

	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

// AddBlock method for BlockChain structure
func (chain *BlockChain) AddBlock(transactions []*Transaction) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.ValueCopy(nil)

		return err
	})
	Handle(err)

	newBlock := CreateBlock(transactions, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})
	Handle(err)
}

// Iterator method for blockchain structure to return
// the original structure to different type of structure
func (chain *BlockChain) Iterator() *BlockChainIterator {
	iterator := &BlockChainIterator{chain.LastHash, chain.Database}

	return iterator
}

// Next method for BlockChainIterator structure
func (iterator *BlockChainIterator) Next() *Block {
	var block *Block
	var encodedBlock []byte

	err := iterator.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iterator.CurrentHash)
		Handle(err)
		err = item.Value(func(val []byte) error {
			encodedBlock = val
			return nil
		})
		block = Deserialize(encodedBlock)

		return err
	})
	Handle(err)

	iterator.CurrentHash = block.PrevHash

	return block
}

// FindUnspentTransactions method for blockchain structure
func (chain *BlockChain) FindUnspentTransactions(address string) []Transaction {
	var unspentTxs []Transaction

	spentTXOs := make(map[string][]int)

	iterator := chain.Iterator()

	for {
		block := iterator.Next()

		for  _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
			Outputs:
				for outIdx, out := range tx.Outputs {
					if spentTXOs[txID] != nil {
						for _, spentOut := range spentTXOs[txID] {
							if spentOut == outIdx {
								continue Outputs
							}
						}
					}
					if out.CanBeUnlocked(address) {
						unspentTxs = append(unspentTxs, *tx)
					}
					if tx.IsCoinbase() == false {
						for _, in := range tx.Inputs {
							if in.CanUnlock(address) {
								inTxID := hex.EncodeToString(in.ID)
								spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
							}
						}
					}
				}
		}

		if len(block.PrevHash) == 0 { // means this block is Genesis block
			break
		}
	}
	return unspentTxs
}

// FindUTXO method
func (chain *BlockChain) FindUTXO(address string) []TxOutput {
	var UTXOs []TxOutput
	unspentTransactions := chain.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}
	return UTXOs
}

// FindSpendableOutputs method for blockchain structure
func (chain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := chain.FindUnspentTransactions(address)
	accumulated := 0

	Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			if out.CanBeUnlocked(address) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				if accumulated > amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOuts

}
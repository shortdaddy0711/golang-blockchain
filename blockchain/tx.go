package blockchain

import (
	"bytes"
	"encoding/gob"

	"github.com/shortdaddy0711/golang-blockchain/wallet"
)

// TxOutput transaction output structure
type TxOutput struct {
	Value      int
	PubKeyHash []byte
}

// TxOutputs structure
type TxOutputs struct {
	Outputs []TxOutput
}

// TxInput transaction input structure
type TxInput struct {
	ID        []byte
	Out       int
	Signature []byte
	PubKey    []byte
}

// NewTXOutput function to give value to TxOutput structure
func NewTXOutput(value int, address string) *TxOutput {
	txo := &TxOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}

// Serialize method for TxOutputs
func (outs TxOutputs) Serialize() []byte {
	var buffer bytes.Buffer
	encode := gob.NewEncoder(&buffer)
	err := encode.Encode(outs)
	Handle(err)
	return buffer.Bytes()
}

// DeserializeOutputs function to deserialize the serialized output
func DeserializeOutputs(data []byte) TxOutputs {
	var outputs TxOutputs
	decode := gob.NewDecoder(bytes.NewReader(data))
	err := decode.Decode(&outputs)
	Handle(err)
	return outputs
}

// UsesKey method for TxInput structure
func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallet.PublicKeyHash(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

// Lock method for TxOutput structure
func (out *TxOutput) Lock(address []byte) {
	pubKeyHash := wallet.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash) - 4]
	out.PubKeyHash = pubKeyHash
}

// IsLockedWithKey method
func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}
package wallet

import (
	"log"

	"github.com/mr-tron/base58"
)

// Base58Encode function that encode input data with Base58 algorithm
func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)

	return []byte(encode)
}

// Base58Decode function that decode the encoded input data
func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]))
	if err != nil {
		log.Panic(err)
	}
	return decode
}






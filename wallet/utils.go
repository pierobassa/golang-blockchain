package wallet

import (
	"github.com/mr-tron/base58"
	"log"
)

/* ------------------- ENCODE & DECODE w/ Base58 ------------------- */
func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)

	return []byte(encode) //returns encode (string) as slice of bytes
}

func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]))

	if err != nil {
		log.Panicln(err)
	}

	return decode
}

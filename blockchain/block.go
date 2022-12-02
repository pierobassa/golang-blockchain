package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int //The nonce is the number that blockchain miners are solving for.
}

/*
*
@returns new pointer to a Block
*/
func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0} //block is a reference (&) to a block created with it's constructor. For now Hash is an empty bytes array

	//Running the Proof of Work algorithm on the block
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

/*
We need a Genesis block due to the fact that each block references to a previous block.
*/
func Genesis() *Block {
	return CreateBlock("Genesis", []byte{}) //Genesis block will have an empty slice of bytes as the previous Hash
}

/* -------------- SERIALIZATION & DESERIALIZATION -------------- */
//BadgerDB only accepts bytes as keys and values. We need a way to serialize and deserialize block data being stored

/*
Serialize() is a method of the Block struct

@returns: slice of bytes representing the serialization of the block
*/
func (b *Block) Serialize() []byte {
	var res bytes.Buffer

	encoder := gob.NewEncoder(&res) //Package gob manages streams of gobs - binary values exchanged between an Encoder (transmitter) and a Decoder (receiver).

	err := encoder.Encode(b)

	if err != nil {
		log.Panic(err)
	}

	return res.Bytes() //returning directly res because we encoded res by passing a reference to res when creating the Encoder
}

/*
@parameters: slice of bytes which represents a block that has been encoded and needs to be decoded
@returns pointer to the Block
*/
func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	if err != nil {
		log.Panic(err)
	}

	return &block
}

/*
Function that handles errors making code cleaner
@parameters: error
*/
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

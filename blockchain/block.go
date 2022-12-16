package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type Block struct {
	Hash         []byte
	Transactions []*Transaction //Array of transactions. A block must contain at least 1 transaction
	PrevHash     []byte
	Nonce        int //The nonce is the number that blockchain miners are solving for.
}

/*
*
@returns new pointer to a Block
*/
func CreateBlock(txs []*Transaction, prevHash []byte) *Block {
	block := &Block{[]byte{}, txs, prevHash, 0} //block is a reference (&) to a block created with it's constructor.

	//Running the Proof of Work algorithm on the block
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

/*
We need a Genesis block due to the fact that each block references to a previous block.
@param 'coinbase' is the first transaction
*/
func Genesis(coinbase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinbase}, []byte{}) //Genesis block will have an empty slice of bytes as the previous Hash
}

/*
Provides a hash for all of the transactions in a block combined.
*/
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte //2D slice of bytes
	var txHash [32]byte   //slice of bytes with a length of 32

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID) //each row of txHashes will be the transaction ID of the current tx in all of the block's transactions
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{})) //putting together each row a a single array of bytes. []bytes{} means that each row will not be separated by anything

	return txHash[:]
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

	Handle(err)

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

	Handle(err)

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

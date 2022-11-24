package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

type Blockchain struct {
	blocks []*Block //Blockchain is an array of pointers to Blocks
}

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
}

/**
	This a function for a Block struct. In particular a function for a reference to a block struct
	It creates the hash for the block and assigns it the the block pointed to by b.
*/
func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{}) //we take a 2D arraymade of Data and PrevHash and we combine it with an empty array of bytes
	hash := sha256.Sum256(info) //we create the hash with SHA256 hashing algorithm

	//assign to the block's hash the hash we obtained
	b.Hash = hash[:] //In go, Arrays and Slices are slightly different and cannot be used interchangeably; however, you can make a slice from an array easily using the [:] operator.
}

/**
	@returns new pointer to a Block
*/
func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash}; //block is a reference (&) to a block created with it's constructor. For now Hash is an empty bytes array
	block.DeriveHash()
	
	return block
}

/*
	This is a method for a Blockchain struct.
	It adds a new block to the blockchain.
*/
func (chain *Blockchain) AddBlock(data string){
	prevBlock := chain.blocks[len(chain.blocks) - 1] //previous block (current last block before adding the new one)
	new := CreateBlock(data, prevBlock.Hash)

	chain.blocks = append(chain.blocks, new); //Adding the new block
}

/*
	We need a Genesis block due to the fact that each block references to a previous block.
*/
func Genesis() *Block {
	return CreateBlock("Genesis", []byte{}) //Genesis block will have an empty slice of bytes as the previous Hash
}

/*
	We initialize our blockchain with the first block which is the Genesis Block
*/
func InitBlockchain() *Blockchain {
	return &Blockchain{[]*Block{Genesis()}} //we create an array of blocks with 1 block which is the genesis block
}

func main() {
	chain := InitBlockchain()

	chain.AddBlock("First Block after Gensis")
	chain.AddBlock("Second Block after Gensis")
	chain.AddBlock("Third Block after Gensis")

	fmt.Printf("Blocks in the chain:\n")
	for _, block := range chain.blocks { //We run a for loop on chain.blocks and each element will be called 'block'
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in Block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("\n")
	}
}
package main

import (
	"fmt"
	"strconv"

	"github.com/pierobassa/golang-blockchain/blockchain"
)

func main() {
	chain := blockchain.InitBlockchain()

	chain.AddBlock("First Block after Gensis")
	chain.AddBlock("Second Block after Gensis")
	chain.AddBlock("Third Block after Gensis")

	fmt.Printf("Blocks in the chain:\n")
	for _, block := range chain.Blocks { //We run a for loop on chain.blocks and each element will be called 'block'
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in Block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("\n")

		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate())) //Proof of work is done on each block, it doesn't store the blocks. Blockchain does 
		fmt.Println()
	}
}
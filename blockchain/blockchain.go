package blockchain

type Blockchain struct {
	Blocks []*Block //Blockchain is an array of pointers to Blocks
}

/*
This is a method for a Blockchain struct.
It adds a new block to the blockchain.
*/
func (chain *Blockchain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1] //previous block (current last block before adding the new one)
	new := CreateBlock(data, prevBlock.Hash)

	chain.Blocks = append(chain.Blocks, new) //Adding the new block
}

/*
We initialize our blockchain with the first block which is the Genesis Block
*/
func InitBlockchain() *Blockchain {
	return &Blockchain{[]*Block{Genesis()}} //we create an array of Blocks with 1 block which is the genesis block
}

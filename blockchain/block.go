package blockchain

type Blockchain struct {
	Blocks []*Block //Blockchain is an array of pointers to Blocks
}

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
This is a method for a Blockchain struct.
It adds a new block to the blockchain.
*/
func (chain *Blockchain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1] //previous block (current last block before adding the new one)
	new := CreateBlock(data, prevBlock.Hash)

	chain.Blocks = append(chain.Blocks, new) //Adding the new block
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
	return &Blockchain{[]*Block{Genesis()}} //we create an array of Blocks with 1 block which is the genesis block
}

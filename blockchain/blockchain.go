package blockchain

import (
	"fmt"
	"github.com/dgraph-io/badger"
)

// Path of the database
const (
	dbPath = "./tmp/blocks"
)

/*
Blockchain struct
  - LastHash: slice of bytes representing the last hash (hash of the last block in the blockchain)
  - Database: pointer to the Badger database
*/
type Blockchain struct {
	Blocks   []*Block //Blockchain is an array of pointers to Blocks
	Database *badger.DB
}

/*
This is a method for a Blockchain struct.
It adds a new block to the blockchain.
*/
func (chain *Blockchain) AddBlock(data string) {
	var lastHash []byte

	opts := badger.DefaultOptions //DefaultOptions sets a list of recommended options for good performance.
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)

	Handle(err)

	//Update allows us to do Read and Write operations. Meanwhile, 'View' is read-only
	err := db.Update(func(txn *badger.Txn) error { //'func(txn *badger.Txn) error' is a closure (anonymous function) which has a pointer to a badger transaction and returns an error if it is occurred
		/*
			1) check if the blockchain has already been stored in the database
			oss the '_' is blank identifier and avoids to declare all returned variables of the function. It allows us to call a function and not have to use one or more return values

			2) If there is a blockchain there then we will create a new blockchain instance in memory and we'll get the last hash of the last block in the blockchain
			The last hash is important because that's how we derive a new block's hash

			3) If there is no existing blockchain in the database we'll create the Genesis block and store it in the database
			 next we'll save the genesis' block's hash as the last block hash (lastHash) in the database
			 then we'll create a new blockchain instance with the lastHash pointing towards the Gensis block
		*/

		//we are getting the value of the last hash which we save last hash in the db with 'lh'
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound { //if the error is equal to ErrKeyNotFound then the blockchain doesn't exist in the database
			fmt.Println("No existing blockchain found")
			genesis := Genesis()
			fmt.Println("Genesis proved")

			err = txn.Set(genesis.Hash, genesis.Serialize()) //We are adding the genesis block to the database with key=genesis' block hash and value=genesis block serialized as slice of bytes
			Handle(err)

			//Setting new data in the database: setting last hash to the genesis block hash
			err = txn.Set([]byte("lh"), genesis.Hash)
			Handle(err)

			lastHash = genesis.Hash

			return err
		} else { //If the database already exists

		}

	})
}

/*
We initialize our blockchain with the first block which is the Genesis Block
*/
func InitBlockchain() *Blockchain {
	return &Blockchain{[]*Block{Genesis()}} //we create an array of Blocks with 1 block which is the genesis block
}

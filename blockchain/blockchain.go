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
	LastHash []byte
	Database *badger.DB
}

/*
Struct used to iterate through the blockchain in the database
*/
type BlockchainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

/*
We initialize our blockchain (if it isn't already present) with the first block which is the Genesis Block
*/
func InitBlockchain() *Blockchain {
	var lastHash []byte

	opts := badger.DefaultOptions(dbPath) //DefaultOptions sets a list of recommended options for good performance.

	db, err := badger.Open(opts)
	Handle(err)

	//Update allows us to do Read and Write operations. Meanwhile, 'View' is read-only
	err = db.Update(func(txn *badger.Txn) error { //'func(txn *badger.Txn) error' is a closure (anonymous function) which has a pointer to a badger transaction and returns an error if it is occurred
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

			lastHash = genesis.Hash

			return err
		} else { //If the database already exists
			item, err := txn.Get([]byte("lh")) //Get the last hash from the db
			Handle(err)

			err = item.Value(func(val []byte) error {
				lastHash = append([]byte{}, val...)

				return nil
			})

			return err
		}
	})

	Handle(err)

	blockchain := Blockchain{lastHash, db}
	return &blockchain
}

/*
This is a method for a Blockchain struct.
It adds a new block to the blockchain.
*/
func (chain *Blockchain) AddBlock(data string) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error { //Read-only transaction to the db
		item, err := txn.Get([]byte("lh"))
		Handle(err)

		err = item.Value(func(val []byte) error {
			lastHash = append([]byte{}, val...)

			return nil
		})

		return err
	})

	Handle(err)

	newBlock := CreateBlock(data, lastHash)

	//Now we need to add the block to the database
	//and update the last hash in the database
	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)

		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})

	Handle(err)
}

/* ------------------ ITERATOR METHODS ------------------- */
/*
Function for the Blockchain struct
Iterator needed to go through the blocks in the blockchain saved in the DB

@returns Pointer to a blochcian iterator which is an iterator for our blockchain
*/
func (chain *Blockchain) Iterator() *BlockchainIterator {
	iter := &BlockchainIterator{chain.LastHash, chain.Database}

	return iter
}

/*
We want to iterate 'backwords'. Which means that we are iterating from the most recent block to the oldest (Genesis block)
@returns a pointer to the next block in the blockchain
*/
func (iter *BlockchainIterator) Next() *Block {
	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)

		var encodedBlock []byte

		err = item.Value(func(val []byte) error { //Since Badger 1.6.0 we need to get the value this way and not with simple 'Value()'
			encodedBlock = append([]byte{}, val...)

			return nil
		})

		//Alternative to Value()  you could also use item.ValueCopy().
		//ex: encodedBlock, err = item.ValueCopy(nil)

		block = Deserialize(encodedBlock)

		return err
	})

	Handle(err)

	iter.CurrentHash = block.PrevHash //We are now chaning the iterator to the previous block

	return block
}

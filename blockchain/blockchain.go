package blockchain

import (
	"encoding/hex"
	"fmt"
	"github.com/dgraph-io/badger"
	"os"
	"runtime"
)

// Path of the database
const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST" //to verify if the blockchain db exists (BadgerDB creates this file on initialization of the DB)
	genesisData = "First Transaction from Genesis"
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

func DBExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

func ContinueBlockchain(address string) *Blockchain {
	if DBExists() == false {
		fmt.Println("No existing blockchain found, create one!")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions(dbPath) //DefaultOptions sets a list of recommended options for good performance.

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)

		err = item.Value(func(val []byte) error {
			lastHash = append([]byte{}, val...)

			return nil
		})

		return err
	})

	Handle(err)

	blockchain := Blockchain{lastHash, db}

	return &blockchain
}

/*
We initialize our blockchain (if it isn't already present) with the first block which is the Genesis Block

@param 'address': address of who inits the blockchain
*/
func InitBlockchain(address string) *Blockchain {
	var lastHash []byte

	//Check if DB already exists
	if DBExists() {
		fmt.Println("Blockchain already exists! No need to init the blockchain again")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(dbPath) //DefaultOptions sets a list of recommended options for good performance.

	db, err := badger.Open(opts)
	Handle(err)

	//Update allows us to do Read and Write operations. Meanwhile, 'View' is read-only
	err = db.Update(func(txn *badger.Txn) error { //'func(txn *badger.Txn) error' is a closure (anonymous function) which has a pointer to a badger transaction and returns an error if it is occurred
		/*
		 We'll create the Genesis block and store it in the database
		 next we'll save the genesis' block's hash as the last block hash (lastHash) in the database
		*/

		//we are getting the value of the last hash which we save last hash in the db with 'lh'
		cbtx := CoinbaseTx(address, genesisData) //Coinbase tx (address is the address that will mine the genesis block and be rewarded the 100 tokens)
		genesis := Genesis(cbtx)

		fmt.Println("Genesis created!")

		//Adding block to the DB
		err = txn.Set(genesis.Hash, genesis.Serialize())
		Handle(err)

		//Setting the last hash in the DB
		err = txn.Set([]byte("lh"), genesis.Hash)

		lastHash = genesis.Hash

		return err
	})

	Handle(err)

	blockchain := Blockchain{lastHash, db}
	return &blockchain
}

/*
This is a method for a Blockchain struct.
It adds a new block to the blockchain.
*/
func (chain *Blockchain) AddBlock(transactions []*Transaction) {
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

	newBlock := CreateBlock(transactions, lastHash)

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

/*
Finds transactions that have outputs NOT referenced by inputs of other transactions.
These are important because if an input hasn't been spent that means that there are tokens that still exist
for a certain user.

# By summing the unspent transactions for a user we can know his balance of tokens

@param: address: the address of the user
*/
func (chain *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTxs []Transaction

	//spent transaction outputs
	//Creating a map where the keys are strings and values are array (slice) of Integers.
	//The array of integers stores the indexes of the outputs SPENT of the transaction (the key is the Transaction's ID)
	spentTXOs := make(map[string][]int)

	iter := chain.Iterator()

	for {
		block := iter.Next()

		//Iterate through the transactions of the current block
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID) //into hexadecimal in string format

		Outputs: //label that helps us continue to this for from the inner for loop (the one for the outputs) and not the external for loop
			for outIdx, out := range tx.Outputs {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx { //This output is a spent output, it can't be part of the UTXs
							continue Outputs //we go to the next output
						}
					}
				}

				//if we reach here, the transaction (tx) is an UNSPENT transaction. Now if the user with address 'address' can unlock the output transaction then it's his and we add it to the array of unspent tranasactions
				if out.CanBeUnlocked(address) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}

			//An output is spent if its index is inside a transaction's input
			if tx.isCoinbase() == false { //a coinbase tx does not have inputs
				for _, in := range tx.Inputs {
					if in.CanUnlock(address) {
						inTxID := hex.EncodeToString(in.ID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out) //Adding the output's index to the key (transaction ID) of the map
					}
				}
			}
		}

		if len(block.PrevHash) == 0 { //if we reach genesis block
			break
		}
	}

	return unspentTxs
}

/*
Finds the Unspent Transaction Outputs of a given address
*/
func (chain *Blockchain) FindUTXO(address string) []TxOutput {
	var UTXOs []TxOutput
	unspentTxs := chain.FindUnspentTransactions(address)

	for _, tx := range unspentTxs {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) { //If the output is unlockable by the address then it is a UTXO of that address
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

/*
This method enables creating normal transactions and not only Coinbase transactions.
To send tokens from one account to another, we need to find the unspent outputs and assure they have enough tokens inside of them.

@param: address -> address we want to check
@param: amount -> the amount we want to send (transfer)

@returns: tuple of int and a map of (string -> array of ints)
*/
func (chain *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := chain.FindUnspentTransactions(address)

	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			if out.CanBeUnlocked(address) && accumulated < amount {
				accumulated += out.Value

				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				if accumulated >= amount {
					break Work //We break from all for loops
				}
			}
		}
	}

	return accumulated, unspentOuts
}

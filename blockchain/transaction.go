package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

/*
Sets the ID for the given Transaction

Creates a hash based on the bytes that represent the transaction
*/
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(tx)
	Handle(err)

	hash = sha256.Sum256(encoded.Bytes())

	tx.ID = hash[:]
}

/*
Coinbase transaction is the first transaction which Input references an empty output because there is no previous transaction
*/
func CoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}

	txin := TxInput{[]byte{}, -1, data} //Empty output, -1 as index because there is no output referenced
	txout := TxOutput{100, to}          //100 tokens in output

	tx := Transaction{nil, []TxInput{txin}, []TxOutput{txout}}
	tx.SetID() //create the hash id for the transaction

	return &tx
}

/*
Creates a new Transaction which is not a Coinbase transaction

@param: from -> from account
@param: to -> to account
@param: amount -> amount of tokens transafered from 'from' to 'to'
@param: chain -> the pointer to the blockchain
*/
func NewTransaction(from string, to string, amount int, chain *Blockchain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	accumulator, validOutputs := chain.FindSpendableOutputs(from, amount)

	if accumulator < amount {
		log.Panic("Error: not enough funds!")
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		Handle(err)

		for _, out := range outs {
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{amount, to})

	if accumulator > amount {
		outputs = append(outputs, TxOutput{accumulator - amount, from}) //Returning the excess amount back to the 'from' account
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}

/*
A transaction is the first transaction (Coinbase) when:
  - The length of the inputs is 1
  - Length of the ID of that input is 0 because we initialize the TxInput with an empty slice of bytes ([]byte{})
  - The input's out index is -1 because there was no previous transaction and the first transaction has the input's out index set to -1
*/
func (tx *Transaction) isCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

type TxOutput struct {
	Value  int    //Number of tokens
	PubKey string //Used to unlock the tokens (in our case the account that made the transaction)
}

type TxInput struct {
	ID  []byte //ID of the Transaction
	Out int    //Position of the Output we are referring to (an Input references an Output)
	Sig string //Used in the Output's PubKey (in our case it is the account that made the transaction)
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
A transaction is the first transaction (Coinbase) when:
  - The length of the inputs is 1
  - Length of the ID of that input is 0 because we initialize the TxInput with an empty slice of bytes ([]byte{})
  - The input's out index is -1 because there was no previous transaction and the first transaction has the input's out index set to -1
*/
func (tx *Transaction) isCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

/* --------------- UNLOCK Data inside the outputs and inputs of a transaction --------------- */

func (in *TxInput) CanUnlock(data string) bool {
	return data == in.Sig //true if the data passed is equal to the signature of the input
}

func (out *TxOutput) CanBeUnlocked(data string) bool {
	return data == out.PubKey //for the output it's the same but this time we are checking if the data is equal to the pub key of the output
}

package blockchain

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

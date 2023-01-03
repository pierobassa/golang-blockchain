package blockchain

type TxOutput struct {
	Value  int    //Number of tokens
	PubKey string //Used to unlock the tokens (in our case the account that made the transaction)
}

type TxInput struct {
	ID  []byte //ID of the Transaction
	Out int    //Position of the Output we are referring to (an Input references an Output)
	Sig string //Used in the Output's PubKey (in our case it is the account that made the transaction)
}

/* --------------- UNLOCK Data inside the outputs and inputs of a transaction --------------- */

func (in *TxInput) CanUnlock(data string) bool {
	return data == in.Sig //true if the data passed is equal to the signature of the input
}

func (out *TxOutput) CanBeUnlocked(data string) bool {
	return data == out.PubKey //for the output it's the same but this time we are checking if the data is equal to the pub key of the output
}

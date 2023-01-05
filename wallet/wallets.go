package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"log"
	"os"
)

/*
We're not using the BadgerDB for storing wallets because we want to use the BadgerDB exclusively for storing the blockchain.
*/

const walletFile = "./tmp/wallets.data"

type Wallets struct {
	Wallets map[string]*Wallet //Map address => Pointer to Wallet
}

/*
Create wallets
*/
func CreateWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFile()

	return &wallets, err
}

/*
Retrieve a Wallet given its address
*/
func (ws Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

/*
Get all Addresses
@returns: array of string
*/
func (ws *Wallets) GetAllAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

/*
Add a new Wallet to all the Wallets

@returns: address (string) of the new Wallet
*/
func (ws *Wallets) AddWallet() string {
	wallet := MakeWallet()

	address := fmt.Sprintf("%s", wallet.Address())

	ws.Wallets[address] = wallet

	return address
}

/*
Saves all the wallets
*/
func (ws *Wallets) SaveFile() {
	var content bytes.Buffer

	gob.Register(elliptic.P256()) //Register the algorithm used. We are registering the type

	encoder := gob.NewEncoder(&content)

	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	err = os.WriteFile(walletFile, content.Bytes(), 0644) //0644 is read & write permissions

	if err != nil {
		log.Panic(err)
	}
}

/*
Loads all the Wallets
*/
func (ws *Wallets) LoadFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	var wallets Wallets

	fileContent, err := os.ReadFile(walletFile)

	if err != nil {
		return err
	}

	gob.Register(elliptic.P256())

	decoder := gob.NewDecoder(bytes.NewReader(fileContent))

	err = decoder.Decode(&wallets)
	if err != nil {
		return err
	}

	ws.Wallets = wallets.Wallets

	return nil //No error returned
}

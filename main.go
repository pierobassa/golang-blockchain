package main

import (
	"github.com/pierobassa/golang-blockchain/cli"
	"os"
)

/*
CLI Commands:
- go run main.go createblockchain -address "John" -> Creates the blockchain with account "John"
- go run main.go printchain -> Prints the blocks in the chain
- go run main.go getbalance -address "John" -> Retrieve the tokens owned by account "John"
- go run .\main.go send -from "John" -to "Fred" -amount 50 -> Send 50 tokens from account "John" to account "Fred"
*/
func main() {
	//In the Go programming language, defer is a keyword that allows a function to postpone the execution of a statement until the surrounding function returns.
	//This can be useful for performing cleanup actions that need to happen regardless of whether the function returns successfully or not.
	defer os.Exit(0) //Further ensures that we properly close the Go application and Badger DB

	cmd := cli.CommandLine{}
	cmd.Run()

	//w := wallet.MakeWallet()
	//w.Address()
}

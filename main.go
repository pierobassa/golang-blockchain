package main

import (
	"flag"
	"fmt"
	"github.com/dgraph-io/badger"
	"os"
	"runtime"
	"strconv"

	"github.com/pierobassa/golang-blockchain/blockchain"
)

/*
CommandLine is a struct to facilitate interacting with the blockchain
*/
type CommandLine struct {
	blockchain *blockchain.Blockchain
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" add -block BLOCK_DATA -> Add a block to the chain")
	fmt.Println(" print -> Prints the blocks in the chain")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()

		runtime.Goexit() //Exits the application but, unlike os.Exit, it exits the application by shutting down the Go routine
		//A goroutine is a lightweight thread managed by the Go runtime.

		//Badger DB has a downside which is it has to garbage collect the values and keys before it shuts down. So if we shut down the application without
		//properly closing the database it can corrupt the data.
	}
}

func (cli *CommandLine) addBlock(data string) {
	cli.blockchain.AddBlock(data)
	fmt.Println("Added the Block!")
}

func (cli *CommandLine) printChain() {
	iter := cli.blockchain.Iterator()

	for {
		block := iter.Next()

		fmt.Println()
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in Block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate())) //Proof of work is done on each block, it doesn't store the blocks. Blockchain does
		fmt.Println()

		//we break out of the for loop when we reach the Genesis block:
		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:]) //parse all arguments from the second one to the last
		blockchain.Handle(err)
	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" { //If the block data was not specified then we will print again usage and exit
			addBlockCmd.Usage()
			runtime.Goexit()
		}

		//dereference the pointer to a string
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

}

/*
CLI Commands:
- go run main.go print -> Prints the blocks in the chain
- go run main.go add -block "[BLOCK DATA]" ->  Add a block to the chain with data: [BLOCK DATA]
*/
func main() {
	//In the Go programming language, defer is a keyword that allows a function to postpone the execution of a statement until the surrounding function returns.
	//This can be useful for performing cleanup actions that need to happen regardless of whether the function returns successfully or not.
	defer os.Exit(0) //Further ensures that we properly close the Go application and Badger DB

	chain := blockchain.InitBlockchain()

	defer func(Database *badger.DB) { //Closes the DB to make sure any pending update is written before closing
		fmt.Println("Closing Badger DB...")
		err := Database.Close()
		blockchain.Handle(err)
	}(chain.Database)

	cli := CommandLine{chain}
	cli.run()
}

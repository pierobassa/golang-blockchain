package main

import (
	"flag"
	"fmt"
	"github.com/dgraph-io/badger"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/pierobassa/golang-blockchain/blockchain"
)

/*
CommandLine is a struct to facilitate interacting with the blockchain
*/
type CommandLine struct {
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" getbalance -address ADDRESS -> get the balance of the ADDRESS")
	fmt.Println(" createblockchain -address ADDRESS -> creates a blockchain")
	fmt.Println(" printchain -> prints the blocks in the chain")
	fmt.Println(" send -from FROM -to TO -amount AMOUNT -> sends AMOUNT from FROM to TO")
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

func (cli *CommandLine) printChain() {
	chain := blockchain.ContinueBlockchain("")

	defer func(Database *badger.DB) { //Closes the DB to make sure any pending update is written before closing
		fmt.Println("Closing Badger DB...")
		err := Database.Close()
		blockchain.Handle(err)
	}(chain.Database)

	iter := chain.Iterator()

	for {
		block := iter.Next()

		fmt.Println()
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
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

func (cli *CommandLine) createBlockchain(address string) {
	chain := blockchain.InitBlockchain(address) //address is the user that mines the genesis block

	err := chain.Database.Close() //Close the DB
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Blockchain created!")
}

func (cli *CommandLine) getBalance(address string) {
	chain := blockchain.ContinueBlockchain(address)

	defer func(Database *badger.DB) { //Closes the DB to make sure any pending update is written before closing
		fmt.Println("Closing Badger DB...")
		err := Database.Close()
		blockchain.Handle(err)
	}(chain.Database)

	balance := 0
	UTXOs := chain.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

func (cli *CommandLine) send(from string, to string, amount int) {
	chain := blockchain.ContinueBlockchain(from)

	defer func(Database *badger.DB) { //Closes the DB to make sure any pending update is written before closing
		fmt.Println("Closing Badger DB...")
		err := Database.Close()
		blockchain.Handle(err)
	}(chain.Database)

	tx := blockchain.NewTransaction(from, to, amount, chain)

	chain.AddBlock([]*blockchain.Transaction{tx})

	fmt.Printf("[SUCCESS SEND] %s -> %d -> %s\n", from, amount, to)
}

func (cli *CommandLine) run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

}

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

	cli := CommandLine{}
	cli.run()
}

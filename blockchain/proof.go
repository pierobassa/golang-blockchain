package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

//Proof of work
/*

- Take the data from the block

- Create a counter (nonce) which start at 0

- Create a hash of the data plus the counter

- Check the hash to see if it meets a set of requirements

REQUIREMENTS:
- First few bytes must cointain a number of zeros

*/

// Here we are hard coding the difficulty but we would also want an algorithm that increases the difficulty after a large period of time
// This is to account for the increasing number of miners on the network and the increasing computation power of the network by all miners.
// We want to make the time to mine a block stay the same and also want to have the block rate stay the same.
const Difficulty = 18

//Proof of the work done for signing a new block:
/*
- Block: pointer to a Block
- Target: pointer to a Big Integer. Integers are accurate but have a limited range. What if you need a really big, accurate number? big.Int (Big integer)
*/
type ProofOfWork struct {
	Block  *Block
	Target *big.Int //Target will be derived from the Difficulty
}

/*
@parameters: pointer to a block
@returns: pointer to a ProofOfWork

We are taking a Block and pairing it with target
*/
func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1) //0000....1

	target.Lsh(target, uint(256-Difficulty)) //we are left shifting target by 256 (bytes of a hash) minus the difficulty

	pow := &ProofOfWork{b, target} //creating the ProofOfWork pointer

	return pow
}

/* -------- PROOFOFWORK STRUCT FUNCTIONS --------- */
/*
 Function on the ProofOfWork struct
 @parameters: nonce(int)
 @returns: slice of bytes

 This function substitues the DeriveHash() function. We are combining Block.PrevHash & Block.Data with an empty slice to create a cohesive set of bytes to return.
*/
func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.HashTransactions(),
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)

	return data
}

/*
We need to create the hash and check if the hash respects the set of requirements.
This main computational function is called Run()

@returns:  int and slice of bytes.
*/
func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0

	fmt.Printf("Finding nonce hash...\n")
	for nonce < math.MaxInt64 { //We are creating a virtual infinite loop
		//1. Prepare data
		data := pow.InitData(nonce)

		//2. Hash into SHA256 format
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)

		//3. Convert the hash into a Big Integer
		intHash.SetBytes(hash[:])

		//4. Compare the Big Integer with the Target Big Integer of ProofOfWork struct
		if intHash.Cmp(pow.Target) == -1 { //Cmp returns -1 means that the hash is less than the target so we have already reached the difficulty needed
			break
		} else {
			nonce++
		}
	}
	fmt.Println()

	return nonce, hash[:]
}

/*
After running the Proof of work run function, we'll have the nonce to derive the hash which must meet the Target that we want
Validate() is used to show that the hash derived is VALID

@returns true if block hash is valid false otherwise
*/
func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)

	hash := sha256.Sum256(data)

	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}

/* -------- UTILITY FUNCTIONS --------- */

/*
Utility function that creates a new bytes buffer to take our number and decode it into bytes.
@parameters: num (64 bit integer)
@returns: slice of bytes
*/
func ToHex(num int64) []byte {
	buff := new(bytes.Buffer) // New buffer for handling bytes

	err := binary.Write(buff, binary.BigEndian, num) //Write() writes the binary representation of num into buff. BigEndian: https://it.wikipedia.org/wiki/Ordine_dei_byte.

	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

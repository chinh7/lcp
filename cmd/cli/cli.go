package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ed25519"

	"github.com/QuoineFinancial/vertex/abi"
	"github.com/QuoineFinancial/vertex/crypto"
)

// type Submission struct {
// 	Method  string   `json:"method"`
// 	Jsonrpc string   `json:"jsonrcp"`
// 	Params  []string `json:"params"`
// 	ID      int      `json:"id"`
// }

func broadcast(serializedTx string) {
	log.Println("Signed Transaction Len:", len(serializedTx))
	msg, _ := hex.DecodeString(serializedTx)
	serializedTx = base64.StdEncoding.EncodeToString(msg)
	log.Println("Params Len:", len(serializedTx))
	fmt.Println(serializedTx)
}

func loadPrivateKey() ed25519.PrivateKey {
	path := os.Getenv("SEED_PATH")
	dat, err := ioutil.ReadFile(path)
	stringData := strings.TrimSuffix(string(dat), "\n")

	parsed, err := hex.DecodeString(stringData)

	if err != nil {
		panic(err)
	}
	return ed25519.NewKeyFromSeed(parsed)
}

func sign(tx *crypto.Tx) {
	privateKey := loadPrivateKey()
	pubkey := make([]byte, ed25519.PublicKeySize)
	copy(pubkey, privateKey[32:])
	tx.From.PubKey = pubkey

	sigHash := tx.GetSigHash()
	signature := ed25519.Sign(privateKey, sigHash)
	tx.From.Signature = signature
}

func deploy(cmd *cobra.Command, args []string) {
	code, err := ioutil.ReadFile(args[1])
	if err != nil {
		panic(err)
	}
	encodedHeader, err := abi.EncodeHeaderToBytes(args[2])
	if err != nil {
		panic(err)
	}

	header, err := abi.DecodeHeader(encodedHeader)
	if err != nil {
		panic(err)
	}

	nonce, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		panic(err)
	}

	data, err := rlp.EncodeToBytes(&abi.Contract{Header: header, Code: code})
	if err != nil {
		panic(err)
	}

	signer := crypto.TxSigner{Nonce: uint64(nonce)}
	tx := &crypto.Tx{Data: data, From: signer}
	sign(tx)
	broadcast(hex.EncodeToString(tx.Serialize()))
}

func invoke(cmd *cobra.Command, args []string) {
	nonce, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		panic(err)
	}
	signer := crypto.TxSigner{Nonce: uint64(nonce)}
	to := crypto.AddressFromString(args[1])

	header, err := abi.LoadHeaderFromFile(args[2])
	if err != nil {
		panic(err)
	}

	function, err := header.GetFunction(args[3])
	if err != nil {
		panic(err)
	}
	encodedArgs, err := abi.EncodeFromString(function.Parameters, args[4:])
	if err != nil {
		panic(err)
	}

	txData := crypto.TxData{Method: args[3], Params: encodedArgs}
	tx := &crypto.Tx{Data: txData.Serialize(), From: signer, To: to}

	sign(tx)
	broadcast(hex.EncodeToString(tx.Serialize()))
}

func main() {
	var cmdDeploy = &cobra.Command{
		Use:   "deploy [nonce] [path to wasm] [path to contract abi json file]",
		Short: "Deploy a wasm contract",
		Args:  cobra.MinimumNArgs(3),
		Run:   deploy,
	}

	var cmdInvoke = &cobra.Command{
		Use:   "invoke [nonce] [address] [path to contract abi json file] [param to invoke]",
		Short: "invoke a smart contract",
		Args:  cobra.MinimumNArgs(4),
		Run:   invoke,
	}

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmdDeploy, cmdInvoke)
	rootCmd.Execute()
}

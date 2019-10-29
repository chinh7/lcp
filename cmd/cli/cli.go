package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

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
	url := "http://localhost:26657/"
	msg, err := hex.DecodeString(serializedTx)
	serializedTx = base64.StdEncoding.EncodeToString(msg)
	log.Println("Params Len:", len(serializedTx))
	// log.Println(serializedTx)
	// if len(serializedTx) > 0 {
	// 	return
	// }
	message := map[string]interface{}{
		"method":  "broadcast_tx_async",
		"id":      123,
		"jsonrpc": "2.0",
		"params":  []string{serializedTx},
	}
	bytesRepresentation, err := json.Marshal(message)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
	defer resp.Body.Close()
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
	data, err := ioutil.ReadFile(args[1])
	if err != nil {
		panic(err)
	}
	encodedHeader, err := abi.EncodeHeaderFromFile(args[2])
	if err != nil {
		panic(err)
	}
	nonce, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		panic(err)
	}
	data = append(encodedHeader, data...)
	signer := crypto.TxSigner{Nonce: uint64(nonce)}
	tx := &crypto.Tx{Data: data, From: signer}
	sign(tx)
	broadcast(hex.EncodeToString(tx.Serialize()))
}

func invoke(cmd *cobra.Command, args []string) {
	var header abi.Header

	nonce, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		panic(err)
	}
	signer := crypto.TxSigner{Nonce: uint64(nonce)}
	to := crypto.AddressFromString(args[1])

	headerFile, err := ioutil.ReadFile(args[2])
	if err != nil {
		panic(err)
	}
	json.Unmarshal(headerFile, &header)

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

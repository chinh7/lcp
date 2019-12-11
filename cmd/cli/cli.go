package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gorilla/rpc/v2/json2"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ed25519"

	"github.com/QuoineFinancial/vertex/abi"
	"github.com/QuoineFinancial/vertex/crypto"
)

// BroadcastRequest of JSON RPC
type BroadcastRequest struct {
	RawTx string `json:"rawTx"`
}

// BroadcastResponse of JSON RPC
type BroadcastResponse struct {
	Hash string `json:"hash"`
}

func broadcast(endpoint, serializedTx string) {
	log.Println("Signed Transaction Len:", len(serializedTx))
	msg, _ := hex.DecodeString(serializedTx)
	serializedTx = base64.StdEncoding.EncodeToString(msg)
	log.Println("Params Len:", len(serializedTx))
	if len(endpoint) > 0 {
		var result BroadcastResponse
		postJSON(endpoint, "chain.Broadcast", BroadcastRequest{RawTx: serializedTx}, &result)
		log.Printf("Transaction hash: %s\n", result.Hash)
	} else {
		log.Println(serializedTx)
	}
}

func loadPrivateKey(path string) ed25519.PrivateKey {
	dat, err := ioutil.ReadFile(path)
	stringData := strings.TrimSuffix(string(dat), "\n")

	parsed, err := hex.DecodeString(stringData)

	if err != nil {
		panic(err)
	}
	return ed25519.NewKeyFromSeed(parsed)
}

func sign(privateKey ed25519.PrivateKey, tx *crypto.Tx) {
	pubkey := make([]byte, ed25519.PublicKeySize)
	copy(pubkey, privateKey[32:])
	tx.From.PubKey = pubkey

	sigHash := tx.GetSigHash()
	signature := ed25519.Sign(privateKey, sigHash)
	tx.From.Signature = signature
}

func deploy(cmd *cobra.Command, args []string) {
	seedPath, endpoint, nonce, gas := parseFlags(cmd)
	privateKey := loadPrivateKey(seedPath)

	code, err := ioutil.ReadFile(args[0])
	if err != nil {
		panic(err)
	}
	encodedHeader, err := abi.EncodeHeaderToBytes(args[1])
	if err != nil {
		panic(err)
	}

	header, err := abi.DecodeHeader(encodedHeader)
	if err != nil {
		panic(err)
	}

	data, err := rlp.EncodeToBytes(&abi.Contract{Header: header, Code: code})
	if err != nil {
		panic(err)
	}

	signer := crypto.TxSigner{Nonce: uint64(nonce)}
	tx := &crypto.Tx{Data: data, From: signer, GasLimit: gas}
	sign(privateKey, tx)
	broadcast(endpoint, hex.EncodeToString(tx.Serialize()))
}

func invoke(cmd *cobra.Command, args []string) {
	seedPath, endpoint, nonce, gas := parseFlags(cmd)
	privateKey := loadPrivateKey(seedPath)

	signer := crypto.TxSigner{Nonce: uint64(nonce)}
	to := crypto.AddressFromString(args[0])

	header, err := abi.LoadHeaderFromFile(args[1])
	if err != nil {
		panic(err)
	}

	function, err := header.GetFunction(args[2])
	if err != nil {
		panic(err)
	}
	encodedArgs, err := abi.EncodeFromString(function.Parameters, args[3:])
	if err != nil {
		panic(err)
	}

	txData := crypto.TxData{Method: args[3], Params: encodedArgs}
	tx := &crypto.Tx{Data: txData.Serialize(), From: signer, To: to, GasLimit: gas}

	sign(privateKey, tx)
	broadcast(endpoint, hex.EncodeToString(tx.Serialize()))
}

func main() {
	var cmdDeploy = &cobra.Command{
		Use:   "deploy [path to wasm] [path to contract abi json file]",
		Short: "Deploy a wasm contract",
		Args:  cobra.MinimumNArgs(2),
		Run:   deploy,
	}

	var cmdInvoke = &cobra.Command{
		Use:   "invoke [address] [path to contract abi json file] [param to invoke]",
		Short: "invoke a smart contract",
		Args:  cobra.MinimumNArgs(3),
		Run:   invoke,
	}

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmdDeploy, cmdInvoke)
	rootCmd.PersistentFlags().StringP("endpoint", "e", "", "Vertex node API endpoint")
	rootCmd.PersistentFlags().Uint64P("gas", "g", 100000, "Gas limit")
	rootCmd.PersistentFlags().StringP("seed", "s", "", "Path to seed")
	rootCmd.PersistentFlags().Uint64P("nonce", "n", 0, "Nonce")
	rootCmd.Execute()
}

func parseFlags(cmd *cobra.Command) (string, string, uint64, uint64) {
	seedPath, err := cmd.Root().Flags().GetString("seed")
	if err != nil {
		panic(err)
	}
	endpoint, err := cmd.Root().Flags().GetString("endpoint")
	if err != nil {
		panic(err)
	}
	nonce, err := cmd.Root().Flags().GetUint64("nonce")
	if err != nil {
		panic(err)
	}
	gas, err := cmd.Root().Flags().GetUint64("gas")
	if err != nil {
		panic(err)
	}
	return seedPath, endpoint, nonce, gas
}

func postJSON(endpoint string, method string, params interface{}, result interface{}) {
	message := map[string]interface{}{
		"method":  method,
		"id":      time.Now().Unix(),
		"jsonrpc": "2.0",
		"params":  params,
	}
	bytesRepresentation, err := json.Marshal(message)
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		panic(err)
	}
	err = json2.DecodeClientResponse(resp.Body, result)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

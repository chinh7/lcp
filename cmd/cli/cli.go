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
	nonce, err := strconv.ParseInt(args[0], 10, 64)
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
	params := make([]interface{}, 0)
	method := args[2]
	for i := 3; i < len(args); i++ {
		param, err := strconv.ParseInt(args[i], 10, 64)
		if err == nil {
			params = append(params, param)
		} else {
			params = append(params, args[i])
		}
	}
	if err != nil {
		panic(err)
	}
	txData := crypto.TxData{Method: method, Params: params}
	tx := &crypto.Tx{Data: txData.Serialize(), From: signer, To: to}

	sign(tx)
	broadcast(hex.EncodeToString(tx.Serialize()))
}

func main() {
	var cmdDeploy = &cobra.Command{
		Use:   "deploy [nonce] [path to wasm]",
		Short: "Deploy a wasm contract",
		Args:  cobra.MinimumNArgs(2),
		Run:   deploy,
	}

	var cmdInvoke = &cobra.Command{
		Use:   "invoke [nonce] [address] [param to invoke]",
		Short: "invoke a smart contract",
		Args:  cobra.MinimumNArgs(3),
		Run:   invoke,
	}

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmdDeploy, cmdInvoke)
	rootCmd.Execute()
}

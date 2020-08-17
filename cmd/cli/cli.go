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

	"github.com/gorilla/rpc/v2/json2"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ed25519"

	"github.com/QuoineFinancial/liquid-chain/api/chain"
	"github.com/QuoineFinancial/liquid-chain/api/storage"
	"github.com/QuoineFinancial/liquid-chain/consensus"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/util"
)

func broadcast(endpoint string, hexTx []byte) {
	log.Println("Signed Transaction Len:", len(hexTx))
	serializedTx := base64.StdEncoding.EncodeToString(hexTx)
	log.Println("Params Len:", len(serializedTx))
	if len(endpoint) > 0 {
		var result chain.BroadcastResult
		postJSON(endpoint, "chain.Broadcast", chain.BroadcastParams{RawTransaction: serializedTx}, &result)

		if result.Code == consensus.ResponseCodeOK {
			log.Println("Broadcast SUCCESS")
			log.Printf("Code: %d\n", result.Code)
			log.Printf("Transaction hash: %s\n", result.TransactionHash)
		} else {
			log.Println("Broadcast FAIL")
			log.Printf("Code: %d\n", result.Code)
			log.Printf("Log: %s\n", result.Log)
		}
	} else {
		log.Println(serializedTx)
	}
}

func loadPrivateKey(path string) ed25519.PrivateKey {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	stringData := strings.TrimSuffix(string(dat), "\n")
	parsed, err := hex.DecodeString(stringData)
	if err != nil {
		panic(err)
	}
	return ed25519.NewKeyFromSeed(parsed)
}

func deploy(cmd *cobra.Command, args []string) {
	seedPath, endpoint, nonce, gas, price, _ := parseFlags(cmd)
	privateKey := loadPrivateKey(seedPath)

	payload, err := util.BuildDeployTxPayload(args[0], args[1], consensus.InitFunctionName, args[2:])
	if err != nil {
		panic(err)
	}

	tx := &crypto.Transaction{
		Version: 1,
		Payload: payload,
		Sender: &crypto.TxSender{
			Nonce:     uint64(nonce),
			PublicKey: privateKey.Public().(ed25519.PublicKey),
		},
		Receiver:  crypto.EmptyAddress,
		GasLimit:  gas,
		GasPrice:  price,
		Signature: nil,
		Receipt:   &crypto.TxReceipt{},
	}
	dataToSign := crypto.GetSigHash(tx)
	tx.Signature = crypto.Sign(privateKey, dataToSign[:])

	if rawTx, err := tx.Encode(); err != nil {
		panic(err)
	} else {
		broadcast(endpoint, rawTx)
	}
}

func invoke(cmd *cobra.Command, args []string) {
	seedPath, endpoint, nonce, gas, price, _ := parseFlags(cmd)
	privateKey := loadPrivateKey(seedPath)

	receiver, err := crypto.AddressFromString(args[0])
	if err != nil {
		panic(err)
	}

	payload, err := util.BuildInvokeTxPayload(args[1], args[2], args[3:])
	if err != nil {
		panic(err)
	}
	tx := &crypto.Transaction{
		Version: 1,
		Payload: payload,
		Sender: &crypto.TxSender{
			Nonce:     uint64(nonce),
			PublicKey: privateKey.Public().(ed25519.PublicKey),
		},
		Receiver:  receiver,
		GasLimit:  gas,
		GasPrice:  price,
		Signature: nil,
		Receipt:   &crypto.TxReceipt{},
	}
	dataToSign := crypto.GetSigHash(tx)
	tx.Signature = crypto.Sign(privateKey, dataToSign[:])

	if rawTx, err := tx.Encode(); err != nil {
		panic(err)
	} else {
		broadcast(endpoint, rawTx)
	}
}

func call(cmd *cobra.Command, args []string) {
	_, endpoint, _, _, _, height := parseFlags(cmd)

	address := args[0]
	method := args[1]
	params := args[2:]

	var result storage.CallResult
	postJSON(endpoint, "storage.Call", storage.CallParams{Address: address, Method: method, Args: params, Height: height}, &result)
	log.Printf("Return: %v", result.Return)
}

func main() {
	var cmdDeploy = &cobra.Command{
		Use:   "deploy [path to wasm] [path to contract abi json file]",
		Short: "Deploy a wasm contract",
		Args:  cobra.MinimumNArgs(2),
		Run:   deploy,
	}

	var cmdInvoke = &cobra.Command{
		Use:   "invoke [address] [path to contract abi json file] [function] [params]",
		Short: "Invoke a smart contract",
		Args:  cobra.MinimumNArgs(3),
		Run:   invoke,
	}

	var cmdCall = &cobra.Command{
		Use:   "call [address] [function] [params]",
		Short: "Call a smart contract (read-only)",
		Args:  cobra.MinimumNArgs(2),
		Run:   call,
	}

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmdDeploy, cmdInvoke, cmdCall)
	rootCmd.PersistentFlags().StringP("endpoint", "e", "", "Vertex node API endpoint")
	rootCmd.PersistentFlags().Uint32P("gas", "g", 100000, "Gas limit")
	rootCmd.PersistentFlags().StringP("seed", "s", "", "Path to seed")
	rootCmd.PersistentFlags().Uint64P("nonce", "n", 0, "Position of transaction")
	rootCmd.PersistentFlags().Int64("height", 0, "Call the method at height")
	rootCmd.PersistentFlags().Uint32P("price", "p", 1, "Gas price")
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func parseFlags(cmd *cobra.Command) (string, string, uint64, uint32, uint32, int64) {
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
	gas, err := cmd.Root().Flags().GetUint32("gas")
	if err != nil {
		panic(err)
	}
	price, err := cmd.Root().Flags().GetUint32("price")
	if err != nil {
		panic(err)
	}
	height, err := cmd.Root().Flags().GetInt64("height")
	if err != nil {
		panic(err)
	}
	return seedPath, endpoint, nonce, gas, price, height
}

func postJSON(endpoint string, method string, params interface{}, result interface{}) {
	message := map[string]interface{}{
		"method":  method,
		"id":      time.Now().Unix(),
		"jsonrpc": "2.0",
		"params":  params,
	}
	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
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

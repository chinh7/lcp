package node

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/ed25519" // This is used in place of crypto/ed25519 to support older version of Go

	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/QuoineFinancial/liquid-chain/api"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/google/go-cmp/cmp"
	"github.com/tendermint/tendermint/config"
)

type testCase struct {
	name   string
	method string
	params string
	result string
}

type testServer struct {
	node *LiquidNode
}

const (
	blockchainTestName = "integration_test"
	gasContractAddress = "LACWIGXH6CZCRRHFSK2F4BINXGUGUS2FSX5GSYG3RMP5T55EV72DHAJ7"
	SEED               = "0c61093a4983f5ba8cf83939efc6719e0c61093a4983f5ba8cf83939efc6719e"
)

func (ts *testServer) startNode() {
	conf := config.ResetTestRoot(blockchainTestName)
	fmt.Println("Init node config data...")

	ts.node = New(conf.RootDir, gasContractAddress)
	conf, err := ts.node.parseConfig()
	if err != nil {
		panic(err)
	}
	conf.LogLevel = "error"
	conf.Consensus.CreateEmptyBlocks = false

	go func() {
		err := ts.node.startTendermintNode(conf)
		if err != nil && err.Error() != "http: Server closed" {
			panic(err)
		}
	}()
	// Wait some time for server to ready
	time.Sleep(4 * time.Second)
}

// Please remember to call stopNode after done testing
func (ts *testServer) stopNode() {
	time.Sleep(2 * time.Second)

	ts.node.stopNode()
	fmt.Println("Clean up node data")
	time.Sleep(500 * time.Millisecond)
	os.RemoveAll(ts.node.rootDir)

	time.Sleep(500 * time.Millisecond)
}

func createDeployTx() string {
	var (
		err           error
		code          []byte
		contractCode  []byte
		encodedHeader []byte
		header        *abi.Header
	)

	if code, err = ioutil.ReadFile("./testdata/contract.wasm"); err != nil {
		panic(err)
	}

	if encodedHeader, err = abi.EncodeHeaderToBytes("./testdata/contract-abi.json"); err != nil {
		panic(err)
	}

	if header, err = abi.DecodeHeader(encodedHeader); err != nil {
		panic(err)
	}

	if contractCode, err = rlp.EncodeToBytes(&abi.Contract{Header: header, Code: code}); err != nil {
		panic(err)
	}

	txData := crypto.TxData{ContractCode: contractCode}
	signer := crypto.TxSigner{Nonce: uint64(0)}
	tx := &crypto.Tx{Data: txData.Serialize(), From: signer, GasLimit: 1, GasPrice: 1}

	privKey := loadPrivateKey(SEED)
	if err = tx.Sign(privKey); err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(tx.Serialize())
}

func createInvokeTx(contractAddress string, nonce uint64, functionName string, params []string) string {
	to, err := crypto.AddressFromString(contractAddress)
	if err != nil {
		panic(err)
	}

	header, err := abi.LoadHeaderFromFile("./testdata/contract-abi.json")
	if err != nil {
		panic(err)
	}

	function, err := header.GetFunction(functionName)
	if err != nil {
		panic(err)
	}

	encodedArgs, err := abi.EncodeFromString(function.Parameters, params)
	if err != nil {
		panic(err)
	}

	signer := crypto.TxSigner{Nonce: uint64(nonce)}
	txData := crypto.TxData{Method: functionName, Params: encodedArgs}
	tx := &crypto.Tx{Data: txData.Serialize(), From: signer, To: to, GasLimit: 1, GasPrice: 1}

	privKey := loadPrivateKey(SEED)
	if err = tx.Sign(privKey); err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(tx.Serialize())
}

func loadPrivateKey(SEED string) ed25519.PrivateKey {
	hexSeed, err := hex.DecodeString(SEED)
	if err != nil {
		panic(err)
	}
	return ed25519.NewKeyFromSeed(hexSeed)
}

func TestBroadcastTx(t *testing.T) {
	ts := &testServer{}
	defer ts.stopNode()
	ts.startNode()

	api := api.NewAPI(":5555", api.Config{
		HomeDir: ts.node.rootDir,
		NodeURL: "tcp://localhost:26657",
		DB:      ts.node.app.StateDB,
	})

	router := api.Router
	testcases := []testCase{
		{
			name:   "Broadcast",
			method: "chain.Broadcast",
			params: fmt.Sprintf(`{"rawTx": "%s"}`, createDeployTx()),
			result: `{"jsonrpc":"2.0","result":{"hash":"53E3715C74FCFCC008AA9E2D7E99C51F109FFCC4EFBFA524D9BA6469EF4F5453","code":0,"log":""},"id":1}`,
		},
	}

	for _, test := range testcases {
		response := httptest.NewRecorder()
		request, _ := makeRequest(test.method, test.params)
		router.ServeHTTP(response, request)
		result := readBody(response)
		if diff := cmp.Diff(string(result), test.result); diff != "" {
			t.Errorf("%s: expect %s, got %s, diff: %s", test.name, test.result, result, diff)
		}
	}
}

func makeRequest(method string, params string) (*http.Request, error) {
	var body string
	if params == "" {
		body = fmt.Sprintf(`{"jsonrpc": "2.0", "id": 1, "method": "%s"}`, method)
	} else {
		body = fmt.Sprintf(`{"jsonrpc": "2.0", "id": 1, "method": "%s", "params": %s}`, method, params)
	}

	req, err := http.NewRequest("POST", "/", bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func readBody(res *httptest.ResponseRecorder) string {
	content, _ := ioutil.ReadAll(res.Body)
	stringResponse := strings.TrimSuffix(string(content), "\n")
	return string(stringResponse)
}

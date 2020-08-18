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

	"github.com/QuoineFinancial/liquid-chain/api"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/util"
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

	ts.node = New(conf.RootDir, "")
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
	})

	router := api.Router
	seed := make([]byte, 32)
	privateKey := ed25519.NewKeyFromSeed(seed)
	sender := crypto.TxSender{
		Nonce:     uint64(0),
		PublicKey: privateKey.Public().(ed25519.PublicKey),
	}
	payload, err := util.BuildDeployTxPayload("./testdata/contract.wasm", "./testdata/contract-abi.json", "init", []string{})
	if err != nil {
		t.Fatal(err)
	}
	deployTx := &crypto.Transaction{
		Sender:    &sender,
		Payload:   payload,
		Receiver:  crypto.EmptyAddress,
		GasLimit:  0,
		GasPrice:  1,
		Signature: nil,
		Receipt:   &crypto.TxReceipt{},
	}
	dataToSign := crypto.GetSigHash(deployTx)
	deployTx.Signature = crypto.Sign(privateKey, dataToSign[:])
	rawTx, _ := deployTx.Encode()
	serializedTx := base64.StdEncoding.EncodeToString(rawTx)
	testcases := []testCase{
		{
			name:   "Broadcast",
			method: "chain.Broadcast",
			params: fmt.Sprintf(`{"rawTx": "%s"}`, serializedTx),
			result: `{"jsonrpc":"2.0","result":{"hash":"7E43B44AC44FFA3FAF53078D6BBCC55ACC7D9BB01AE20860D617226F69A168F5","code":0,"log":""},"id":1}`,
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

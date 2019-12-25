package node

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/QuoineFinancial/liquid-chain/api"
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
)

func (ts *testServer) startNode() error {
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
	time.Sleep(2 * time.Second)

	return nil
}

// Please remember to call stopNode after done testing
func (ts *testServer) stopNode() {
	ts.node.stopNode()
	fmt.Println("Clean up node data")
	err := os.RemoveAll(ts.node.rootDir)
	if err != nil {
		panic(err)
	}
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

	contractHex, err := ioutil.ReadFile("./testdata/contract-hex.txt")
	if err != nil {
		panic(err)
	}
	body := fmt.Sprintf(`{"rawTx": "%s"}`, string(contractHex))
	testcases := []testCase{
		{"Broadcast", "chain.Broadcast", body, `{"jsonrpc":"2.0","result":{"hash":"A1D8A3AB7CC2971CD26409418205D1DAEEF5AEBB228BFA8643ABA3AEE126B961"},"id":1}`},
		{"Broadcast", "chain.Broadcast", `{"rawTx": "+KH4ZKBZheHjjJtrb75IOm2u18eUHwn1+rEw8fI5kPueGVO7V4C4QA40TwFk+Muh6/5vUsM0szRyaW8g0iWrALLj+DdebvFSWcgbXEeR7m1WUxGIz+W/Wy3N3ka668fzE6gXNLM6tQGR0IRtaW50ismIhAMAAAAAAACjWAVkGufwsijE5ZK0XgUNuahqS0WV+mlg24sf2fekr/QzgT+DAYagAQ=="}`, `{"jsonrpc":"2.0","result":{"hash":"F058970C18F36659C6722A7BD6656E01AB425158B553B58BE6AD79F54025FC63"},"id":1}`},
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

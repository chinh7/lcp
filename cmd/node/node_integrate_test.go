package node

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
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

func (ts *testServer) startNode() {
	conf := config.ResetTestRoot(blockchainTestName)
	fmt.Println("Init node config data...")

	ts.node = NewNode(conf.RootDir, gasContractAddress)
	conf, err := ts.node.parseConfig()
	if err != nil {
		panic(err)
	}
	conf.LogLevel = "error"
	conf.Consensus.CreateEmptyBlocks = false

	go func() {
		err := ts.node.startNode(conf, true)
		if err != nil && err.Error() != "http: Server closed" {
			panic(err)
		}
	}()

	// Wait some time for server to ready
	time.Sleep(5 * time.Second)
}

// Please remember to call stopNode after done testing
func (ts *testServer) stopNode() {
	time.Sleep(2 * time.Second)

	ts.node.stopNode()
	fmt.Println("Clean up node data")
	err := os.RemoveAll(ts.node.rootDir)
	if err != nil {
		panic(err)
	}

	// Wait for blockchain to completely close
	time.Sleep(500 * time.Millisecond)
}

func TestBroadcastTx(t *testing.T) {
	ts := &testServer{}
	defer ts.stopNode()
	ts.startNode()

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
		result := MakeRequest(test.method, test.params)
		if result != test.result {
			t.Errorf("%s: expect %s, got %s", test.name, test.result, result)
		}
	}
}

func MakeRequest(method string, params string) string {
	client := resty.New()
	var body string
	if params == "" {
		body = fmt.Sprintf(`{"jsonrpc": "2.0", "id": 1, "method": "%s"}`, method)
	} else {
		body = fmt.Sprintf(`{"jsonrpc": "2.0", "id": 1, "method": "%s", "params": %s}`, method, params)
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post("http://localhost:5555")

	if err != nil {
		panic(err)
	}

	return resp.String()
}

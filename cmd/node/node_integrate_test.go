package node

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
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

func setupNode() (*LiquidNode, *config.Config, error) {
	conf := config.ResetTestRoot("integration_test")
	fmt.Printf("Init node data in %s\n", conf.RootDir)

	gasContractAddress := "LACWIGXH6CZCRRHFSK2F4BINXGUGUS2FSX5GSYG3RMP5T55EV72DHAJ7"
	liquidNode := NewNode(conf.RootDir, gasContractAddress)

	conf, err := liquidNode.parseConfig()
	conf.LogLevel = "error"
	conf.Consensus.CreateEmptyBlocks = false

	return liquidNode, conf, err
}

func startNode(wg *sync.WaitGroup, liquidNode *LiquidNode, conf *config.Config) {
	defer wg.Done()
	err := liquidNode.startNode(conf, true)
	if err != nil && err.Error() != "http: Server closed" {
		panic(err)
	}
}

func TestBroadcastTx(t *testing.T) {
	liquidNode, conf, err := setupNode()
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	wg.Add(2)

	contractHex, err := ioutil.ReadFile("./testdata/contract-hex.txt")
	if err != nil {
		panic(err)
	}
	body := fmt.Sprintf(`{"rawTx": "%s"}`, string(contractHex))
	testcases := []testCase{
		{"Broadcast", "chain.Broadcast", body, `{"jsonrpc":"2.0","result":{"hash":"A1D8A3AB7CC2971CD26409418205D1DAEEF5AEBB228BFA8643ABA3AEE126B961"},"id":1}`},
		{"Broadcast", "chain.Broadcast", `{"rawTx": "+KH4ZKBZheHjjJtrb75IOm2u18eUHwn1+rEw8fI5kPueGVO7V4C4QA40TwFk+Muh6/5vUsM0szRyaW8g0iWrALLj+DdebvFSWcgbXEeR7m1WUxGIz+W/Wy3N3ka668fzE6gXNLM6tQGR0IRtaW50ismIhAMAAAAAAACjWAVkGufwsijE5ZK0XgUNuahqS0WV+mlg24sf2fekr/QzgT+DAYagAQ=="}`, `{"jsonrpc":"2.0","result":{"hash":"F058970C18F36659C6722A7BD6656E01AB425158B553B58BE6AD79F54025FC63"},"id":1}`},
	}

	go startNode(&wg, liquidNode, conf)
	go func() {
		defer wg.Done()
		// Wait for blockchain node and api server to start
		time.Sleep(5 * time.Second)

		for _, test := range testcases {
			result := MakeRequest(test.method, test.params)
			if result != test.result {
				t.Errorf("%s: expect %s, got %s", test.name, test.result, result)
			}
		}

		time.Sleep(3 * time.Second)

		liquidNode.stopNode(true)
		fmt.Printf("Removing data in %s\n", liquidNode.rootDir)
		err := os.RemoveAll(liquidNode.rootDir)
		if err != nil {
			panic(err)
		}

		// Wait for blockchain to completely close
		time.Sleep(2 * time.Second)
	}()

	wg.Wait()
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

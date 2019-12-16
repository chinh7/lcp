package node

import (
	"fmt"
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
	conf.LogLevel = "*:error"
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

	testcases := []testCase{
		{"Broadcast", "chain.Broadcast", `{"rawTx": "+Nn4ZKBZheHjjJtrb75IOm2u18eUHwn1+rEw8fI5kPueGVO7VwO4QLMiPEIy9zhnECC8tqMB3f1XpJbicg8y9f1Rci9nZrvDAUs8jQ3yrzMTFVjt+kfQzmfuKZ2128DP1e4n/zk/AA64Sm1pbnQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADJiOgDAAAAAAAAo1gFZBrn8LIoxOWStF4FDbmoaktFlfppYNuLH9n3pK/0M4E/gicQ"}`, `{"jsonrpc":"2.0","result":{"hash":"6B5E4F16974D1D387BFA41B1BF606A107D6361831ED60D3E1CCE5A37C7490DA8"},"id":1}`},
	}

	go startNode(&wg, liquidNode, conf)
	go func() {
		defer wg.Done()
		// Wait for blockchain node and api server to start
		time.Sleep(5 * time.Second)

		for _, test := range testcases {
			result := MakeRequest(test.method, test.params)
			if result != test.result {
				t.Errorf("%s: expected %s, got %s", test.name, test.result, result)
			}
		}

		liquidNode.StopNode(true)
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

package resource

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/tendermint/tendermint/cmd/tendermint/commands"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/lite/proxy"
	"github.com/tendermint/tendermint/rpc/client"
)

// TendermintAPI is client to interact with Tendermint RPC
type TendermintAPI = client.Client

func readChainID(homeDir string) string {
	genesisPath := filepath.Join(homeDir, "/config/genesis.json")
	configFile, err := os.Open(genesisPath)
	if err != nil {
		panic("Unable to read genesis.json with error:\n" + err.Error())
	}
	configBytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		panic("Invalid format of genesis.json")
	}
	var config struct {
		ChainID string `json:"chain_id"`
	}
	if err := json.Unmarshal(configBytes, &config); err != nil {
		panic("Could not read chain_id from genesis file")
	}
	return config.ChainID
}

// NewTendermintAPI returns new instance of TendermintAPI
func NewTendermintAPI(homeDir, nodeURL string) TendermintAPI {
	chainID := readChainID(homeDir)
	logFileName := fmt.Sprintf("tendermint-api-%d.log", time.Now().Unix())
	logFilePath := filepath.Join(homeDir, logFileName)
	tendermintLoggerFile, _ := os.Create(logFilePath)
	defer tendermintLoggerFile.Close()
	logger := log.NewTMLogger(log.NewSyncWriter(tendermintLoggerFile))

	cacheSize := 10
	nodeURL, err := commands.EnsureAddrHasSchemeOrDefaultToTCP(nodeURL)
	if err != nil {
		common.Exit(err.Error())
	}
	node := client.NewHTTP(nodeURL, "/websocket")
	cert, err := proxy.NewVerifier(chainID, homeDir, node, logger, cacheSize)
	if err != nil {
		common.Exit(err.Error())
	}
	cert.SetLogger(logger)
	return proxy.SecureClient(node, cert)
}

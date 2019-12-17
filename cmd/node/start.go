package node

import (
	"fmt"
	"os"

	"github.com/QuoineFinancial/liquid-chain/api"
	"github.com/QuoineFinancial/liquid-chain/consensus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/cmd/tendermint/commands"
	"github.com/tendermint/tendermint/config"
	tmflags "github.com/tendermint/tendermint/libs/cli/flags"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	tmNode "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
)

func (node *LiquidNode) newTendermintNode(config *config.Config, logger log.Logger) (*tmNode.Node, error) {
	node.app = consensus.NewApp(config.Moniker, config.DBDir(), node.gasContractAddress)
	nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
	if err != nil {
		return nil, err
	}

	oldPrivVal := config.OldPrivValidatorFile()
	newPrivValKey := config.PrivValidatorKeyFile()
	newPrivValState := config.PrivValidatorStateFile()

	if _, err := os.Stat(oldPrivVal); !os.IsNotExist(err) {
		oldPV, err := privval.LoadOldFilePV(oldPrivVal)
		if err != nil {
			return nil, fmt.Errorf("Error reading OldPrivValidator from %v: %v", oldPrivVal, err)
		}
		logger.Info("Upgrading PrivValidator file",
			"old", oldPrivVal,
			"newKey", newPrivValKey,
			"newState", newPrivValState,
		)
		oldPV.Upgrade(newPrivValKey, newPrivValState)
	}
	return tmNode.NewNode(config,
		privval.LoadOrGenFilePV(newPrivValKey, newPrivValState),
		nodeKey,
		proxy.NewLocalClientCreator(node.app),
		tmNode.DefaultGenesisDocProviderFunc(config),
		tmNode.DefaultDBProvider,
		tmNode.DefaultMetricsProvider(config.Instrumentation),
		logger.With("module", "node"),
	)
}

func (node *LiquidNode) parseConfig() (*config.Config, error) {
	conf := config.DefaultConfig()
	err := viper.Unmarshal(conf)
	if err != nil {
		return nil, err
	}

	conf.SetRoot(node.rootDir)
	config.EnsureRoot(node.rootDir)
	if err = conf.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("error in config file: %v", err)
	}

	return conf, err
}

func (node *LiquidNode) startNode(conf *config.Config, apiFlag bool) error {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	// Set log level by --log_level flag or default
	logger, err := tmflags.ParseLogLevel(conf.LogLevel, logger, config.DefaultLogLevel())
	if err != nil {
		return err
	}

	n, err := node.newTendermintNode(conf, logger)
	if err != nil {
		return fmt.Errorf("Failed to create node: %v", err)
	}
	node.tmNode = n

	// Stop upon receiving SIGTERM or CTRL-C.
	common.TrapSignal(logger, func() {
		if n.IsRunning() {
			_ = n.Stop() // TODO: Properly handle error
		}
	})

	if err := n.Start(); err != nil {
		return fmt.Errorf("Failed to start node: %v", err)
	}
	logger.Info("Started node", "nodeInfo", n.Switch().NodeInfo())

	if apiFlag {
		node.vertexApi = api.NewAPI(":5555", api.Config{
			HomeDir: node.rootDir,
			NodeURL: "tcp://localhost:26657",
			DB:      node.app.StateDB,
		})
		err := node.vertexApi.Serve()
		if err != nil {
			return err
		}
	}

	return nil
}

func (node *LiquidNode) stopNode(apiFlag bool) {
	if node.tmNode.IsRunning() {
		_ = node.tmNode.Stop() // TODO: Properly handle error
	}

	node.vertexApi.Close()
}

func (node *LiquidNode) addStartNodeCommand() {
	var apiFlag bool
	cmd := &cobra.Command{
		Use:   "start [--api]",
		Short: "Start the liquid node",
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := node.parseConfig()
			if err != nil {
				return fmt.Errorf("Failed to parse config: %v", err)
			}

			err = node.startNode(conf, apiFlag)
			if err != nil {
				return err
			}

			// Run forever.
			select {}
		},
	}
	cmd.PersistentFlags().BoolVarP(&apiFlag, "api", "a", false, "start api")

	commands.AddNodeFlags(cmd)
	node.command.AddCommand(cmd)
}

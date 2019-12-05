package node

import (
	"fmt"
	"os"

	"github.com/QuoineFinancial/vertex/api"
	"github.com/QuoineFinancial/vertex/consensus"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/cmd/tendermint/commands"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	tmNode "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
)

func (node *VertexNode) newTendermintNode(logger log.Logger) (*tmNode.Node, error) {
	node.app = consensus.NewApp(node.config.Moniker, node.config.DBDir(), node.gasContractAddress)
	nodeKey, err := p2p.LoadOrGenNodeKey(node.config.NodeKeyFile())
	if err != nil {
		return nil, err
	}

	oldPrivVal := node.config.OldPrivValidatorFile()
	newPrivValKey := node.config.PrivValidatorKeyFile()
	newPrivValState := node.config.PrivValidatorStateFile()

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
	return tmNode.NewNode(&node.config,
		privval.LoadOrGenFilePV(newPrivValKey, newPrivValState),
		nodeKey,
		proxy.NewLocalClientCreator(node.app),
		tmNode.DefaultGenesisDocProviderFunc(&node.config),
		tmNode.DefaultDBProvider,
		tmNode.DefaultMetricsProvider(node.config.Instrumentation),
		logger.With("module", "node"),
	)
}

func (node *VertexNode) addStartNodeCommand() {
	var apiFlag bool
	cmd := &cobra.Command{
		Use:   "start [--api]",
		Short: "Start the vertex node",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
			n, err := node.newTendermintNode(logger)
			if err != nil {
				return fmt.Errorf("Failed to create node: %v", err)
			}

			// Stop upon receiving SIGTERM or CTRL-C.
			common.TrapSignal(logger, func() {
				if n.IsRunning() {
					n.Stop()
				}
			})

			if err := n.Start(); err != nil {
				return fmt.Errorf("Failed to start node: %v", err)
			}
			logger.Info("Started node", "nodeInfo", n.Switch().NodeInfo())

			if apiFlag == true {
				apiServer := api.NewAPI(":5555", api.Config{
					HomeDir: node.config.RootDir,
					NodeURL: "tcp://localhost:26657",
				})
				apiServer.Serve()
			}

			// Run forever.
			select {}
		},
	}
	cmd.PersistentFlags().BoolVarP(&apiFlag, "api", "a", false, "start api")
	commands.AddNodeFlags(cmd)

	node.command.AddCommand(cmd)
}

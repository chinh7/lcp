package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/QuoineFinancial/vertex/api"
	"github.com/QuoineFinancial/vertex/consensus"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/cmd/tendermint/commands"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
)

// Ref: github.com/tendermint/tendermint/cmd/tendermint/main.go
func main() {
	rootCmd := commands.RootCmd
	rootCmd.AddCommand(
		commands.GenValidatorCmd,
		commands.InitFilesCmd,
		commands.ProbeUpnpCmd,
		commands.LiteCmd,
		commands.ReplayCmd,
		commands.ReplayConsoleCmd,
		commands.ResetAllCmd,
		commands.ResetPrivValidatorCmd,
		commands.ShowValidatorCmd,
		commands.TestnetFilesCmd,
		commands.ShowNodeIDCmd,
		commands.GenNodeKeyCmd,
		commands.VersionCmd)
	// NOTE:
	// Users wishing to:
	//	* Use an external signer for their validators
	//	* Supply an in-proc abci app
	//	* Supply a genesis doc file from another source
	//	* Provide their own DB implementation
	// can copy this file and use something other than the
	// DefaultNewNode function
	// Create & start node
	runNodeCmd := commands.NewRunNodeCmd(newNode)
	rootCmd.AddCommand(addAPI(runNodeCmd))
	command := cli.PrepareBaseCmd(rootCmd, "TM", os.ExpandEnv(filepath.Join("$HOME", config.DefaultTendermintDir)))
	if err := command.Execute(); err != nil {
		panic(err)
	}
}

func addAPI(command *cobra.Command) *cobra.Command {
	var apiFlag bool
	newCommand := &cobra.Command{
		Use:   "start [--api]",
		Short: "Start the tendermint node",
		RunE: func(cmd *cobra.Command, args []string) error {
			if apiFlag == true {
				go (func() {
					time.Sleep(time.Second)
					apiServer := api.NewAPI(":5555", api.Config{
						HomeDir: os.ExpandEnv(filepath.Join("$HOME", config.DefaultTendermintDir)),
						NodeURL: "tcp://localhost:26657",
					})
					apiServer.Serve()
				})()
			}
			return command.RunE(cmd, args)
		},
	}
	newCommand.PersistentFlags().BoolVarP(&apiFlag, "api", "a", false, "start api")
	return newCommand
}

// Ref: github.com/tendermint/tendermint/node/node.go (func DefaultNewNode)
func newNode(config *config.Config, logger log.Logger) (*node.Node, error) {
	app := consensus.NewApp("sample app")
	// Generate node PrivKey
	nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
	if err != nil {
		return nil, err
	}
	// Convert old PrivValidator if it exists.
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
	return node.NewNode(config,
		privval.LoadOrGenFilePV(newPrivValKey, newPrivValState),
		nodeKey,
		proxy.NewLocalClientCreator(app),
		node.DefaultGenesisDocProviderFunc(config),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(config.Instrumentation),
		logger.With("module", "node"),
	)
}

package node

import (
	"github.com/QuoineFinancial/vertex/consensus"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/cmd/tendermint/commands"
	"github.com/tendermint/tendermint/libs/cli"
)

// VertexNode is the space where app and command lives
type VertexNode struct {
	rootDir            string
	gasContractAddress string
	app                *consensus.App
	command            *cobra.Command
}

// New returns new instance of Node
func New(rootDir string, gasContractAddress string) *VertexNode {
	vertexNode := VertexNode{
		rootDir:            rootDir,
		command:            commands.RootCmd,
		gasContractAddress: gasContractAddress,
	}
	vertexNode.addDefaultCommands()
	vertexNode.addStartNodeCommand()
	return &vertexNode
}

func (node *VertexNode) addDefaultCommands() {
	node.command.AddCommand(
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
		commands.VersionCmd,
	)

}

// Execute run the node.command base on user input
func (node *VertexNode) Execute() {
	prefix := "TM"
	command := cli.PrepareBaseCmd(node.command, prefix, node.rootDir)
	if err := command.Execute(); err != nil {
		panic(err)
	}
}

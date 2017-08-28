package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/tendermint/basecoin/cmd/commands"
	"github.com/tendermint/tmlibs/cli"

	// import _ to register the mint plugin to apptx
	_ "github.com/tendermint/basecoin-examples/mintcoin/cmd/mintcoin/commands"
)

func main() {

	var RootCmd = &cobra.Command{
		Use: "mintcoin",
	}

	RootCmd.AddCommand(
		commands.InitCmd,
		commands.StartCmd,
		commands.TxCmd,
		commands.QueryCmd,
		commands.KeyCmd,
		commands.VerifyCmd,
		commands.BlockCmd,
		commands.AccountCmd,
		commands.UnsafeResetAllCmd,
		commands.QuickVersionCmd("0.2.0"),
	)

	cmd := cli.PrepareMainCmd(RootCmd, "MT", os.ExpandEnv("$HOME/.mintcoin"))
	cmd.Execute()
}

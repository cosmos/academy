package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/tendermint/basecoin/cmd/commands"
	"github.com/tendermint/tmlibs/cli"

	// import _ to register escrow and options to apptx
	_ "github.com/tendermint/basecoin-examples/trader/cmd/trader/commands"
)

func main() {

	var RootCmd = &cobra.Command{
		Use: "trader",
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

	cmd := cli.PrepareMainCmd(RootCmd, "TR", os.ExpandEnv("$HOME/.trader"))
	cmd.Execute()
}

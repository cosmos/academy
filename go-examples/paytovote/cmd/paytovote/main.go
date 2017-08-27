package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/tendermint/basecoin/cmd/commands"
	"github.com/tendermint/tmlibs/cli"

	_ "github.com/tendermint/basecoin-examples/paytovote/cmd/paytovote/commands"
)

func main() {

	var RootCmd = &cobra.Command{
		Use: "paytovote",
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

	cmd := cli.PrepareMainCmd(RootCmd, "PV", os.ExpandEnv("$HOME/.paytovote"))
	cmd.Execute()
}

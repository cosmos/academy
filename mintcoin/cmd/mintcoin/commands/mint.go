package commands

import (
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/tendermint/basecoin-examples/mintcoin"
	bcmd "github.com/tendermint/basecoin/cmd/commands"
	"github.com/tendermint/basecoin/types"
	wire "github.com/tendermint/go-wire"
)

const MintName = "mint"

var (
	//flags
	MintToFlag     string
	MintAmountFlag string

	//Commands
	MintTxCmd = &cobra.Command{
		Use:   "mint",
		Short: "Craft a transaction to mint some more currency",
		RunE:  mintTxCmd,
	}
)

func init() {

	//register flags
	flags := []bcmd.Flag2Register{
		{&MintToFlag, "mintto", "", "Where to send the newly minted coins"},
		{&MintAmountFlag, "mint", "", "Amount of coins to mint in format <amt><coin>,<amt2><coin2>,..."},
	}
	bcmd.RegisterFlags(MintTxCmd, flags)

	bcmd.RegisterTxSubcommand(MintTxCmd)
	bcmd.RegisterStartPlugin(MintName, func() types.Plugin { return mintcoin.New(MintName) })
}

func mintTxCmd(cmd *cobra.Command, args []string) error {

	// convert destination address to bytes
	to, err := hex.DecodeString(bcmd.StripHex(MintToFlag))
	if err != nil {
		return errors.Errorf("To address is invalid hex: %v\n", err)
	}

	amountCoins, err := types.ParseCoins(MintAmountFlag)
	if err != nil {
		return err
	}

	mintTx := mintcoin.MintTx{
		Credits: []mintcoin.Credit{
			{
				Addr:   to,
				Amount: amountCoins,
			},
		},
	}
	fmt.Println("MintTx:", string(wire.JSONBytes(mintTx)))
	data := wire.BinaryBytes(mintTx)

	return bcmd.AppTx(MintName, data)
}

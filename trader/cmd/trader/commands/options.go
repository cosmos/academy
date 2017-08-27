package commands

import (
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	bcmd "github.com/tendermint/basecoin/cmd/commands"
	bc "github.com/tendermint/basecoin/types"
	wire "github.com/tendermint/go-wire"

	"github.com/tendermint/basecoin-examples/trader/plugins/options"
	"github.com/tendermint/basecoin-examples/trader/types"
)

const OptionName = "options"

var (
	//flags
	OptionNodeFlag        string
	OptionAddrFlag        string
	OptionExpireFlag      uint64
	OptionSellToFlag      string
	OptionTradeAmountFlag string
	OptionPriceAmountFlag string

	//commands
	CmdOptionsTx = &cobra.Command{
		Use:   "options",
		Short: "Create, trade, and exercise currency options",
	}

	CmdOptionsCreateTx = &cobra.Command{
		Use:   "create",
		Short: "Create a new option by sending money",
		RunE:  cmdOptionCreateTx,
	}

	CmdOptionsSellTx = &cobra.Command{
		Use:   "sell",
		Short: "Offer to sell this option",
		RunE:  cmdOptionSellTx,
	}

	CmdOptionsBuyTx = &cobra.Command{
		Use:   "buy",
		Short: "Attempt to buy this option",
		RunE:  cmdOptionBuyTx,
	}

	CmdOptionsExerciseTx = &cobra.Command{
		Use:   "exercise",
		Short: "Exercise this option to trade currency at the given rate",
		RunE:  cmdOptionExerciseTx,
	}

	CmdOptionsDissolveTx = &cobra.Command{
		Use:   "disolve",
		Short: "Attempt to disolve this option (if never sold, or already expired)",
		RunE:  cmdOptionDissolveTx,
	}

	CmdOptionsQuery = &cobra.Command{
		Use:   "query [address]",
		Short: "Return the contents of the given option",
		RunE:  cmdOptionQuery,
	}
)

func init() {

	//Register Flags
	createTxFlags := []bcmd.Flag2Register{
		{&OptionExpireFlag, "expire", uint64(0), "The block height when the option expires"},
		{&OptionTradeAmountFlag, "trade", "", "Amount of coins to trade in format <amt><coin>,<amt2><coin2>,..."},
	}
	addrFlag := bcmd.Flag2Register{
		&OptionAddrFlag, "option", "", "The address of this option"}
	sellTxFlags := []bcmd.Flag2Register{
		addrFlag,
		{&OptionSellToFlag, "sellto", "", "Who to sell the options to (optional)"},
		{&OptionPriceAmountFlag, "price", "", "Price to buy option in format <amt><coin>,<amt2><coin2>,..."},
	}
	buyTxFlags := []bcmd.Flag2Register{
		addrFlag,
	}
	exerciseTxFlags := []bcmd.Flag2Register{
		addrFlag,
	}
	dissolveTxFlags := []bcmd.Flag2Register{
		addrFlag,
	}
	queryFlags := []bcmd.Flag2Register{
		{&OptionNodeFlag, "node", "tcp://localhost:46657", "Tendermint RPC address"},
	}
	bcmd.RegisterFlags(CmdOptionsCreateTx, createTxFlags)
	bcmd.RegisterFlags(CmdOptionsSellTx, sellTxFlags)
	bcmd.RegisterFlags(CmdOptionsBuyTx, buyTxFlags)
	bcmd.RegisterFlags(CmdOptionsExerciseTx, exerciseTxFlags)
	bcmd.RegisterFlags(CmdOptionsDissolveTx, dissolveTxFlags)
	bcmd.RegisterFlags(CmdOptionsQuery, queryFlags)

	//Register Subcommands
	CmdOptionsTx.AddCommand(
		CmdOptionsCreateTx,
		CmdOptionsSellTx,
		CmdOptionsBuyTx,
		CmdOptionsExerciseTx,
		CmdOptionsDissolveTx,
		CmdOptionsQuery,
	)

	bcmd.RegisterTxSubcommand(CmdOptionsTx)
	bcmd.RegisterStartPlugin(OptionName,
		func() bc.Plugin { return options.New(OptionName) })
}

func cmdOptionCreateTx(cmd *cobra.Command, args []string) error {

	tradeCoins, err := bc.ParseCoins(OptionTradeAmountFlag)
	if err != nil {
		return err
	}

	tx := types.CreateOptionTx{
		Expiration: OptionExpireFlag,
		Trade:      tradeCoins,
	}
	data := types.OptionsTxBytes(tx)
	return bcmd.AppTx(OptionName, data)
}

func cmdOptionSellTx(cmd *cobra.Command, args []string) error {

	// convert destination address to bytes
	addr, err := hex.DecodeString(bcmd.StripHex(OptionAddrFlag))
	if err != nil {
		return errors.Errorf("Recv address is invalid hex: %v\n", err)
	}

	buyer, err := hex.DecodeString(bcmd.StripHex(OptionSellToFlag))
	if err != nil { // this is optional, we can ignore it
		buyer = nil
	}

	priceCoins, err := bc.ParseCoins(OptionPriceAmountFlag)
	if err != nil {
		return err
	}

	tx := types.SellOptionTx{
		Addr:      addr,
		NewHolder: buyer,
		Price:     priceCoins,
	}
	data := types.OptionsTxBytes(tx)
	return bcmd.AppTx(OptionName, data)
}

func cmdOptionBuyTx(cmd *cobra.Command, args []string) error {

	// convert destination address to bytes
	addr, err := hex.DecodeString(bcmd.StripHex(OptionAddrFlag))
	if err != nil {
		return errors.Errorf("Recv address is invalid hex: %v\n", err)
	}

	tx := types.BuyOptionTx{
		Addr: addr,
	}
	data := types.OptionsTxBytes(tx)
	return bcmd.AppTx(OptionName, data)
}

func cmdOptionExerciseTx(cmd *cobra.Command, args []string) error {

	// convert destination address to bytes
	addr, err := hex.DecodeString(bcmd.StripHex(OptionAddrFlag))
	if err != nil {
		return errors.Errorf("Recv address is invalid hex: %v\n", err)
	}

	tx := types.ExerciseOptionTx{
		Addr: addr,
	}
	data := types.OptionsTxBytes(tx)
	return bcmd.AppTx(OptionName, data)
}

func cmdOptionDissolveTx(cmd *cobra.Command, args []string) error {

	// convert destination address to bytes
	addr, err := hex.DecodeString(bcmd.StripHex(OptionAddrFlag))
	if err != nil {
		return errors.Errorf("Recv address is invalid hex: %v\n", err)
	}

	tx := types.DisolveOptionTx{
		Addr: addr,
	}
	data := types.OptionsTxBytes(tx)
	return bcmd.AppTx(OptionName, data)
}

func cmdOptionQuery(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("account command requires an argument ([address])") //never stack trace
	}
	addrHex := bcmd.StripHex(args[0])

	// convert destination address to bytes
	addr, err := hex.DecodeString(addrHex)
	if err != nil {
		return errors.Errorf("Recv address is invalid hex: %v\n", err)
	}

	opt, err := getOption(OptionNodeFlag, addr)
	if err != nil {
		return err
	}
	fmt.Println(string(wire.JSONBytes(opt)))
	return nil
}

func getOption(tmAddr string, address []byte) (*types.OptionData, error) {
	prefix := []byte(fmt.Sprintf("%s/", OptionName))
	key := append(prefix, address...)
	response, err := bcmd.Query(tmAddr, key)
	if err != nil {
		return nil, err
	}

	optionBytes := response.Value

	if len(optionBytes) == 0 {
		return nil, fmt.Errorf("Option bytes are empty for address: %X ", address)
	}
	opt, err := types.ParseOptionData(optionBytes)
	if err != nil {
		return nil, fmt.Errorf("Error reading option %X error: %v",
			optionBytes, err.Error())
	}
	return &opt, nil
}

package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/tharsis/evmos/x/intrarelayer/types"
)

// NewTxCmd returns a root CLI command handler for certain modules/ibc-evm transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "ibc-evm tx subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// txCmd.AddCommand()
	return txCmd
}

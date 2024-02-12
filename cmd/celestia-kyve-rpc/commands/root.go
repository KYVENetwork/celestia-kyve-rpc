package commands

import (
	"fmt"
	"github.com/KYVENetwork/celestia-kyve-rpc/utils"
	"github.com/spf13/cobra"
)

var (
	logger = utils.CelestiaKyveRpcLogger("commands")
)

var (
	chainId      string
	port         int64
	restEndpoint string
	storageRest  string
)

// RootCmd is the root command for celestia-kyve-rpc.
var rootCmd = &cobra.Command{
	Use:   "celestia-kyve-rpc",
	Short: "The first trustless Celestia RPC, providing validated data through KYVE.",
}

func Execute() {
	versionCmd.Flags().SortFlags = false

	if err := rootCmd.Execute(); err != nil {
		panic(fmt.Errorf("failed to execute root command: %w", err))
	}
}

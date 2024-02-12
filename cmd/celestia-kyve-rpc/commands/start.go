package commands

import (
	"fmt"
	"github.com/KYVENetwork/celestia-kyve-rpc/server"
	"github.com/KYVENetwork/celestia-kyve-rpc/utils"
	"github.com/spf13/cobra"
	"strings"
)

func init() {
	startCmd.Flags().StringVar(&chainId, "chain-id", utils.DefaultChainId, fmt.Sprintf("KYVE chain id [\"%s\",\"%s\", \"%s\"]", utils.ChainIdMainnet, utils.ChainIdKaon, utils.ChainIdKorellia))

	startCmd.Flags().Int64Var(&port, "port", 4242, "API server port")

	startCmd.Flags().StringVar(&restEndpoint, "rest-endpoint", "", "KYVE API endpoint to retrieve validated bundles")

	startCmd.Flags().StringVar(&storageRest, "storage-rest", "", "storage endpoint for requesting bundle data")

	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Celestia-KYVE-RPC",
	Run: func(cmd *cobra.Command, args []string) {
		endpoint := utils.GetChainRest(chainId, restEndpoint)
		storageRest = strings.TrimSuffix(storageRest, "/")

		server.StartApiServer(chainId, endpoint, storageRest, port)
	},
}

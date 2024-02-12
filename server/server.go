package server

import (
	"fmt"
	"github.com/KYVENetwork/celestia-kyve-rpc/collectors/bundles"
	"github.com/KYVENetwork/celestia-kyve-rpc/types"
	"github.com/KYVENetwork/celestia-kyve-rpc/utils"
	"github.com/gin-gonic/gin"
	"github.com/tendermint/tendermint/libs/json"
	"net/http"
	"strconv"
)

var (
	logger = utils.CelestiaKyveRpcLogger("server")
)

type ApiServer struct {
	chainId      string
	restEndpoint string
	storageRest  string
}

// TODO: Replace with Source-Registry integration
var (
	MainnetPoolMap  = make(map[string]int64)
	KaonPoolMap     = make(map[string]int64)
	KorelliaPoolMap = make(map[string]int64)
)

func StartApiServer(chainId, restEndpoint, storageRest string, port string) *ApiServer {
	apiServer := &ApiServer{
		chainId:      chainId,
		restEndpoint: restEndpoint,
		storageRest:  storageRest,
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.GET("/GetAll", apiServer.GetAll)

	if err := r.Run(fmt.Sprintf(":%v", port)); err != nil {
		logger.Error().Str("err", err.Error()).Msg("failed to run api server")
	}

	return apiServer
}

func (apiServer *ApiServer) GetAll(c *gin.Context) {
	heightStr := c.Query("height")
	namespace := c.Query("namespace")

	// TODO: Replace with Source-Registry integration
	KorelliaPoolMap["AAAAAAAAAAAAAAAAAAAAAAAAAIZiad33fbxA7Z0="] = 73

	var poolId int64

	switch apiServer.chainId {
	case utils.ChainIdMainnet:
		poolId = MainnetPoolMap[namespace]
	case utils.ChainIdKaon:
		poolId = KaonPoolMap[namespace]
	case utils.ChainIdKorellia:
		poolId = KorelliaPoolMap[namespace]
	}

	height, err := strconv.Atoi(heightStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	compressedBundle, err := bundles.GetBundleByKey(height, apiServer.restEndpoint, poolId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	decompressedBundle, err := bundles.GetDataFromFinalizedBundle(*compressedBundle, apiServer.storageRest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("failed to decompress bundle: %v", err.Error()),
		})
		return
	}

	// parse bundle
	var bundle types.Bundle

	if err := json.Unmarshal(decompressedBundle, &bundle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("failed to unmarshal bundle data: %v", err.Error()),
		})
		return
	}

	for _, dataItem := range bundle {
		itemHeight, err := strconv.Atoi(dataItem.Key)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("failed to parse block height from key: %v", err.Error()),
			})
			return
		}

		// skip blocks until we reach start height
		if itemHeight < height {
			continue
		} else if itemHeight == height {
			c.JSON(http.StatusOK, dataItem.Value)
			return
		}
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": fmt.Sprintf("failed to find data item in bundle"),
	})
	return
}

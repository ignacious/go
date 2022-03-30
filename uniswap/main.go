package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/machinebox/graphql"
)

type tokenResp struct {
	ID     string   `json:"id"`
	Volume float64  `json:"volume"`
	Pools  []string `json:"pools"`
}

type blockResp struct {
	Number int      `json:"id"`
	Swaps  []string `json:"swaps"`
	Assets []string `json:"assets"`
}

func main() {
	router := gin.Default()
	router.GET("/assets/:id", getAssetByID)
	router.GET("/blocks/:number", getBlockByNumber)
	router.Run("localhost:8080")
}

func getBlockByNumber(c *gin.Context) {
	number := c.Param("number")

	client := graphql.NewClient("https://api.thegraph.com/subgraphs/name/ianlapham/uniswap-v3-alt")
	blockNumber, _ := strconv.Atoi(number)

	query := graphql.NewRequest(`
		query ($number: Int!) {
			swaps(block:{number:$number}) {
				id,
				token0 {
				id
				},
				token1 {
				id
				}
			}
		}
	`)

	query.Var("number", int(blockNumber))

	ctx := context.Background()
	var responseData interface{}

	if err := client.Run(ctx, query, &responseData); err != nil {
		panic(err)
	}

	m := responseData.(map[string]interface{})

	swaps := m["swaps"]

	var swapSlice []string
	var tokenSlice []string

	for _, v := range swaps.([]interface{}) {
		swap := v.(map[string]interface{})

		if swapId, ok := swap["id"].(string); ok {
			swapSlice = append(swapSlice, swapId)
		}

		token0 := swap["token0"]
		token1 := swap["token1"]

		if token0Id, ok := token0.(map[string]interface{})["id"].(string); ok {
			tokenSlice = append(tokenSlice, token0Id)
		}

		if token1Id, ok := token1.(map[string]interface{})["id"].(string); ok {
			tokenSlice = append(tokenSlice, token1Id)
		}
	}

	swapSlice = RemoveDuplicatesFromSlice(swapSlice)
	tokenSlice = RemoveDuplicatesFromSlice(tokenSlice)

	var resp = blockResp{
		Number: blockNumber,
		Swaps:  swapSlice,
		Assets: tokenSlice,
	}

	c.IndentedJSON(http.StatusOK, resp)
}

func getAssetByID(c *gin.Context) {
	assetId := c.Param("id")

	//Default to extreme timestamps. Better way do this but this essentially includes all whitelisted pools.
	poolCreatedAtStart := c.DefaultQuery("poolCreatedAtStart", "0")
	poolCreatedAtEnd := c.DefaultQuery("poolCreatedAtEnd", "9999999999")

	fmt.Println("start:", poolCreatedAtStart, "end:", poolCreatedAtEnd)
	client := graphql.NewClient("https://api.thegraph.com/subgraphs/name/ianlapham/uniswap-v3-alt")
	query := graphql.NewRequest(`
	query ($assetId: String!, $poolCreatedAtStart: String!, $poolCreatedAtEnd: String!) {
		token (id:$assetId) {
			whitelistPools(where:
			{
			  createdAtTimestamp_gte: $poolCreatedAtStart, 
				createdAtTimestamp_lte: $poolCreatedAtEnd
			}) {
			id, createdAtBlockNumber, createdAtTimestamp,
			token0{id}, token1{id},
			volumeToken0,
			volumeToken1
		  }
		},
	  }
    `)

	query.Var("assetId", assetId)
	query.Var("poolCreatedAtStart", poolCreatedAtStart)
	query.Var("poolCreatedAtEnd", poolCreatedAtEnd)

	ctx := context.Background()
	var responseData interface{}

	if err := client.Run(ctx, query, &responseData); err != nil {
		panic(err)
	}

	m := responseData.(map[string]interface{})
	token := m["token"]

	whiteListPools := token.(map[string]interface{})
	pools := whiteListPools["whitelistPools"]

	tokenVolume := float64(0)
	var slice []string

	for _, v := range pools.([]interface{}) {
		pool := v.(map[string]interface{})
		poolId := pool["id"]
		if idStr, ok := poolId.(string); ok {
			slice = append(slice, idStr)
		}

		token0 := pool["token0"]
		token1 := pool["token1"]

		token0Id := token0.(map[string]interface{})["id"]
		token1Id := token1.(map[string]interface{})["id"]

		if token0Id == assetId {
			vt := v.(map[string]interface{})["volumeToken0"]
			if str, ok := vt.(string); ok {
				volumeToken, _ := strconv.ParseFloat(str, 32)
				tokenVolume = tokenVolume + volumeToken
			}
		}

		if token1Id == assetId {
			vt := v.(map[string]interface{})["volumeToken1"]
			if str, ok := vt.(string); ok {
				volumeToken, _ := strconv.ParseFloat(str, 32)
				tokenVolume = tokenVolume + volumeToken
			}
		}
	}
	var resp = tokenResp{
		ID:     assetId,
		Volume: tokenVolume,
		Pools:  slice,
	}

	c.IndentedJSON(http.StatusOK, resp)
}

//https://siongui.github.io/2018/04/14/go-remove-duplicates-from-slice-or-array/
func RemoveDuplicatesFromSlice(s []string) []string {
	m := make(map[string]bool)
	for _, item := range s {
		if _, ok := m[item]; ok {
			// duplicate item
			fmt.Println(item, "is a duplicate")
		} else {
			m[item] = true
		}
	}

	var result []string
	for item, _ := range m {
		result = append(result, item)
	}
	return result
}

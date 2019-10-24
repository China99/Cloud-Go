package handler

import (
	"Cloud-Go/handler/eth"
	"Cloud-Go/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	NodeUrl = "https://ropsten.infura.io/v3/bd40ce91ed6b4c26a6e8b219eca567db"
	Version = "latest"
)

//单次：查询余额
func GetBalance(c *gin.Context) {
	//获取address
	address := c.Request.FormValue("address")

	requester := eth.NewETHRPCRequester(NodeUrl)
	balance, err := requester.GetETHBalanceV(address, Version)
	//转化为实际金额

	eth2float64, err := util.Eth2float64(balance)
	if err != nil {
		//查询失败
		err := "地址查询失败"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "success",
			"data": eth2float64,
		})
	}
}

//多条查询：查询余额
func GetBalances(c *gin.Context) {
	//获取地址
	address := c.QueryArray("address")
	/*	for _, value := range address {
		fmt.Println("--",value)
	}*/
	requester := eth.NewETHRPCRequester(NodeUrl)
	balances, err := requester.GetBalances(address, Version)
	eth2float64, err := util.Eth2float64(balances)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "",
			"msg":  "查询失败",
			"data": err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "success",
			"data": balances,
		})
	}

}

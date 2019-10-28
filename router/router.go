package router

import (
	"Cloud-Go/handler"
	"github.com/gin-gonic/gin"
)

//Router:路由表配置
func Router() *gin.Engine {
	//gin framework
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	//处理静态资源
	router.Static("/static", "./static")

	//新建钱包
	router.POST("/user/createethwallet", handler.CreateWallet)
	//删除钱包
	router.DELETE("/user/eth/deletewallet", handler.DeleteWallet)
	//查看余额
	router.GET("/user/eth/getbalance", handler.GetBalance)
	//批量查询余额
	router.GET("/user/eth/getbalances", handler.GetBalances)
	//转账
	router.POST("/user/eth/sendeth", handler.SendETHTransaction)
	//根据交易hash查询交易数据
	router.GET("/user/eth/txhash", handler.GetBlockInfoByHash)
	router.GET("/user/eth/txhashs", handler.GetBlockTransactions)
	//根据区块号查询区块信息
	router.GET("/user/eth/getfullblock", handler.GetFullBlockInfo)
	//获取最新区块信息
	router.GET("/user/eth/getlastblocknumber", handler.GetLastBlockNumber)

	//转账思路
	//1.验证身份
	//2.与自己密匙文件进行解密获得转账权限
	//3..验证转账密码与密匙进行比对完成转账，否则报错

	return router
}

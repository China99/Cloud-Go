package handler

import (
	"Cloud-Go/handler/eth"
	"Cloud-Go/util"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//测试变量
const (
	NodeUrl  = "https://ropsten.infura.io/v3/bd40ce91ed6b4c26a6e8b219eca567db"
	NodeUrl1 = "https://mainnet.infura.io/v3/bd40ce91ed6b4c26a6e8b219eca567db"
	Version  = "latest"
)

//新建钱包
func CreateWallet(c *gin.Context) {
	password := c.Query("password")
	//先判断
	if len(password) < 6 {
		fmt.Println(errors.New("创建钱包失败,密码要大于6位！"))

		c.JSON(http.StatusBadRequest, gin.H{
			"code": 0,
			"msg":  "error",
			"data": "创建钱包失败,密码要大于6位！",
		})
		return
	} else {

		address, err := eth.NewETHRPCRequester(NodeUrl).CreateETHWallet(password) // 创建成功

		if err != nil {
			err.Error()
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "success",
			"data": address,
		})
	}

}

//删除钱包
func DeleteWallet(c *gin.Context) {
	//删除钱包只是删除此地址在数据库中的记录
	//区块中的信息无法删除
}

//单次：查询余额
func GetBalance(c *gin.Context) {
	//获取address
	address := c.Request.FormValue("address")

	requester := eth.NewETHRPCRequester(NodeUrl)
	balance, err := requester.GetETHBalanceService(address, Version)
	//转化为实际金额
	var str1 []string
	strings := append(str1, balance)
	eth2float64, err := util.Eth2float64(strings)
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
	balances, err := requester.GetBalancesService(address, Version)
	eth2float64, err := util.Eth2float64(balances)
	fmt.Println(eth2float64)
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
			"data": eth2float64,
		})
	}

}

//转账交易
func SendETHTransaction(c *gin.Context) {

	//from 转账地址
	//to  接收地址
	//amount 金额
	//Gas 消耗
	//加密密码

	/*gasLimit := uint64(100000)
	  gasPrice := uint64(36000000000)*/
	const (
		//用于加密的盐值(自定义)
		pwdSalt = "*#890"
	)
	str1 := "13456aa"
	pwd := util.Sha1([]byte(str1 + pwdSalt))
	fmt.Println(pwd)
	//
	from, _ := c.GetPostForm("from")
	to, _ := c.GetPostForm("to")
	value, _ := c.GetPostForm("value")
	gaslimit, _ := c.GetPostForm("gaslimit")
	gasprice, _ := c.GetPostForm("gasprice")
	password, _ := c.GetPostForm("password")
	//from := c.Query("from")
	//to := c.Query("to")
	//value := c.Query("value")
	//gaslimit := c.Request.FormValue("gaslimit")
	//gasprice := c.Query("gasprice")
	//password := c.Query("password")

	gaslimitUint, _ := strconv.ParseUint(gaslimit, 10, 64)
	gasPriceUint, _ := strconv.ParseUint(gasprice, 10, 64)

	if from == "" || len(from) != 42 {
		_ = errors.New("非法交易address值")
		return
	}
	// 当前这笔交易消耗的油费最大值是 (gasLimit * gasPrice) / 10^18 ETH
	//password md5加密以后，通过md5解码
	//pwd2:=util.Sha1([]byte(password + pwdSalt))
	if pwd == password {
		err := util.UnlockETHWallet("./keystores", from, "13456aa")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		txHash, err := eth.NewETHRPCRequester(NodeUrl).SendETHTransactionService(from, to, value, gaslimitUint, gasPriceUint)

		if err != nil {
			//转账失败
			fmt.Println(err.Error())

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg":  "success",
			"code": 1,
			"data": txHash,
		})

	} else {
		fmt.Println("密码错误")
		return
	}

}

//根据区块hash查询交易数据
func GetBlockInfoByHash(c *gin.Context) {

	//txHash := "0x53c5b03e392d6aa68a0df26b6d466ae8fbd1c2c5b74f9baae05434ec9a18a282" //测试区块
	//txHash:="0xf91f63d87fe4599aa673225d459e88a4ce0927f6cb225722cacedf7176d56849"
	txHash := c.Query("txhash")
	if txHash == "" || len(txHash) != 66 {
		// 这里演示，在调用 rpc 接口函数的时候，都要先进行入参的合法性判断
		fmt.Println("非法交易 hash 值")
		return
	}
	txInfo, err := eth.NewETHRPCRequester(NodeUrl1).GetBlockInfoByHashService(txHash)
	if err != nil {
		// 查询失败，打印出信息
		fmt.Println("查询交易失败，信息是：", err.Error())
		return
	}
	// 查询成功，将 transaction 结果的结构体以 json 格式序列化，再以 string 格式输出

	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "success",
		"data": txInfo,
	})
	//	fmt.Println(string(json))
}

//根据区块hash批量获取对应交易信息
//建议单条查询，如果有一个地址对，一个地址不对，返回错误其正确的无法返回
func GetBlockTransactions(c *gin.Context) {

	txHashs := c.QueryArray("txhash")

	if txHashs == nil || len(txHashs) == 0 {
		// 这里演示，在调用 rpc 接口函数的时候，都要先进行入参的合法性判断
		fmt.Println("非法交易 hash 数组")
		return
	}
	txInfos, err := eth.NewETHRPCRequester(NodeUrl).GetTransactionsService(txHashs)
	if err != nil {
		// 查询失败，打印出信息
		fmt.Println("查询交易失败，信息是：", err.Error())
		return
	}
	// 查询成功，将 transaction 结果的结构体以 json 格式序列化，再以 string 格式输出

	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "success",
		"data": txInfos,
	})
}

//根据区块号查询交易信息
func GetFullBlockInfo(c *gin.Context) {

	requester := eth.NewETHRPCRequester(NodeUrl)

	// 获取区块号
	number, _ := requester.GetLatestBlockNumberService()

	// 获取区块信息
	fullBlock, err := requester.GetBlockInfoByNumberService(number)

	if err != nil {
		// 查询失败，打印出信息
		fmt.Println("获取区块信息失败，信息是：", err.Error())
		return
	}
	// 查询成功，将 区块 结果的结构体以 json 格式序列化，再以 string 格式输出

	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "success",
		"data": fullBlock,
	})
	//fmt.Println("根据区块号获取区块信息", string(json1))

	// 根据区块 hash 获取区块信息
	fullBlock, err = requester.GetBlockInfoByHashService(fullBlock.ParentHash)

	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "success",
		"data": fullBlock,
	})
}

//获取最新的区块号
//最好放入定时任务，定时扫描一次
func GetLastBlockNumber(c *gin.Context) {

	number, err := eth.NewETHRPCRequester(NodeUrl).GetLatestBlockNumberService()
	if err != nil {
		// 查询失败，打印出信息
		fmt.Println("获取区块号失败，信息是：", err.Error())
		return
	}
	//num1, _ := json.Marshal(number)
	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "success",
		"data": number,
	})

}

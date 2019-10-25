package eth

import (
	"Cloud-Go/model"
	"Cloud-Go/util"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
)

type ETHRPCRequester struct {
	nonceManager *NonceManager
	client       *EthRpcClient
}

func NewETHRPCRequester(nodeUrl string) *ETHRPCRequester {
	requester := &ETHRPCRequester{}
	// 实例化 nonce 管理器
	requester.nonceManager = NewNonceManager()
	// 实例化 rpc 客户端
	requester.client = NewETHRPCClient(nodeUrl)
	return requester
}

//创建钱包
// 创建以太坊钱包  可以通过
func (r *ETHRPCRequester) CreateETHWallet(password string) (string, error) {

	//路径可变
	keydir := "./keystores" // 它是用来存储所创建了的钱包的 keystore 文件的目录
	// StandardScryptN 是Scrypt加密算法的标准N参数
	// StandardScryptP 是Scrypt加密算法的标准P参数
	ks := keystore.NewKeyStore(keydir, keystore.StandardScryptN, keystore.StandardScryptP)
	wallet, err := ks.NewAccount(password) // 传入密码，创建钱包
	if err != nil {
		return "0x", err
	}
	return wallet.Address.String(), nil
}

// 单条查询：根据以太坊地址，查询以太坊 eth 的余额
func (r *ETHRPCRequester) GetETHBalanceService(address string, version string) (string, error) {
	name := "eth_getBalance"
	result := ""
	// 对应文档，第一个参数就是要被查询的以太坊地址
	// 第二个参数就是 latest
	err := r.client.GetRpc().Call(&result, name, address, version)
	if err != nil {
		return "", err
	}
	if result == "" {
		return "", errors.New("eth balance is null")
	}
	// 因为查询所返回的结果是一个16进制的字符串，
	// 为了方便阅读，我们在下面使用 go 的大数处理将其转为 10 进制，
	// 并防止数位溢出

	ten, _ := new(big.Int).SetString(result[2:], 16)

	return ten.String(), nil

}

//批量查询：根据以太坊地址数组，查询以太坊 eth的余额
func (r *ETHRPCRequester) GetBalancesService(address []string, version string) ([]string, error) {
	name := "eth_getBalance"

	//数组存储每个请求的指针
	var rets []*string
	//获取address数组的长度，方便在循环中诸葛实例化
	size := len(address)
	var reqs []rpc.BatchElem
	for i := 0; i < size; i++ {
		ret := ""
		// 实例化每个 BatchElem
		req := rpc.BatchElem{
			Method: name,
			Args:   []interface{}{address[i], "latest"},
			// &ret 传入单个请求的结果引用，这样是保证它在函数内部被修改值后，回到函数外来，值仍有效
			Result: &ret,
		}
		reqs = append(reqs, req)  // 将每个 BatchElem 添加到 BatchElem 数组
		rets = append(rets, &ret) // 每个请求的结果引用添加到结果数组中
	}

	err := r.client.GetRpc().BatchCall(reqs) //传入数组，发起批量请求
	if err != nil {
		return nil, err
	}
	//查询每个请求是否有错误
	for _, req := range reqs {
		if req.Error != nil {
			return nil, req.Error
		}
	}
	var finalRet []string
	for _, item := range rets {
		ten, _ := new(big.Int).SetString((*item)[2:], 16) //去掉前2位，转成16进制转10进制
		finalRet = append(finalRet, ten.String())
	}
	return finalRet, err

}

// 根据以太坊地址获取 nonce
func (n *NonceManager) GetNonceService(address string) *big.Int {
	if n.nonceMemCache == nil {
		n.nonceMemCache = map[string]*big.Int{}
	}
	n.lock.Lock()         // 加锁
	defer n.lock.Unlock() // 当该函数执行完毕，进行解锁
	return n.nonceMemCache[address]
}

// 获取地址的 nonce 值
func (r *ETHRPCRequester) GetNonce(address string) (uint64, error) {
	methodName := "eth_getTransactionCount"
	nonce := ""
	// 因为我们要查询最新的，根据基于 etTransactionCount 情况下的区块号关系，选取 pending
	err := r.client.client.Call(&nonce, methodName, address, "pending")
	if err != nil {
		return 0, fmt.Errorf("发送交易失败! %s", err.Error())
	}
	n, _ := new(big.Int).SetString(nonce[2:], 16) // 采用大数类型将 16 进制的结果转为 10 进制
	return n.Uint64(), nil                        // 返回交易 hash
}

//以太坊转战交易，根据transaction的不同变量设置，达到不同种类的交易
//返回交易哈希
func (r *ETHRPCRequester) SendTransactionService(address string, transaction *types.Transaction) (string, error) {
	//签名交易数据
	signTx, err := util.SignETHTransaction(address, transaction)
	if err != nil {
		return "", fmt.Errorf("签名失败！%s", err.Error())
	}
	txRlpData, err := rlp.EncodeToBytes(signTx)
	if err != nil {
		return "", fmt.Errorf("rlp序列化失败！%s", err.Error())
	}
	//调用以太坊的rpc接口
	txHash := ""
	methodName := "eth_sendRawTransaction"
	err = r.client.client.Call(&txHash, methodName, hexutil.Encode(txRlpData))
	if err != nil {
		return "", fmt.Errorf("发送交易失败！%s", err.Error())
	}
	return txHash, nil
}

// 发送 ETH 交易，或称转账 ETH
// 参数分别是：交易发起地址，交易接收地址，ETH数量，油费设置
func (r *ETHRPCRequester) SendETHTransactionService(fromStr, toStr, valueStr string, gasLimit, gasPrice uint64) (string, error) {

	if !common.IsHexAddress(fromStr) || !common.IsHexAddress(toStr) {
		return "", errors.New("invalid address")
	}

	to := common.HexToAddress(toStr) // 将字符串类型的转为 address 类型的
	gasPrice_ := new(big.Int).SetUint64(gasPrice)

	// value 乘上 10^decimal，得出真实的转账值，ETH 的小数位是 18 位
	realV := util.GetRealDecimalValue(valueStr, 18)
	if realV == "" {
		return "", errors.New("invalid value")
	}
	amount, _ := new(big.Int).SetString(realV, 10)

	// 获取 nonce
	nonce := r.nonceManager.GetNonceService(fromStr)
	if nonce == nil {
		// nonce 不存在，开始访问节点获取
		n, err := r.GetNonce(fromStr)
		if err != nil {
			return "", fmt.Errorf("获取 nonce 失败 %s", err.Error())
		}
		nonce = new(big.Int).SetUint64(n)
	}

	// 构建 data，因为 eth 交易转账类型，data 是空，我们设置空字符串即可
	data := []byte("")

	// 构建交易结构体
	transaction := types.NewTransaction(
		nonce.Uint64(),
		to,
		amount,
		gasLimit,
		gasPrice_,
		data)

	return r.SendTransactionService(fromStr, transaction)
}

// 根据区块 hash，获取区块信息
func (r *ETHRPCRequester) GetBlockInfoByHashService(blockHash string) (*model.FullBlock, error) {
	methodName := "eth_getBlockByHash"
	fullBlock := model.FullBlock{}
	// eth_getBlockByHash 的第二个参数：
	// 如果是 true，则返回完整的区块信息，false 则 transaction 部分只返回交易hash数组
	err := r.client.client.Call(&fullBlock, methodName, blockHash, true)
	if err != nil {
		return nil, fmt.Errorf("get block info failed! %s", err.Error())
	}
	if fullBlock.Number == "" {
		return nil, fmt.Errorf("block info is empty %s", blockHash)
	}
	return &fullBlock, nil
}

// 根据交易 hash 字符串数组批量获取对应的交易信息
func (r *ETHRPCRequester) GetTransactionsService(txHashs []string) ([]*model.Transaction, error) {
	name := "eth_getTransactionByHash"

	// 结果数组存储的是每个请求的结果指针，也就是引用
	rets := []*model.Transaction{}

	// 获取 hash 数组的长度，方便在循环中逐个实例化 BatchElem
	size := len(txHashs)

	reqs := []rpc.BatchElem{}
	for i := 0; i < size; i++ {
		ret := model.Transaction{}
		// 实例化每个 BatchElem
		req := rpc.BatchElem{
			Method: name,
			Args:   []interface{}{txHashs[i]},
			// &ret 传入单个请求的结果引用，这样是保证它在函数内部被修改值后，回到函数外来，值仍有效
			Result: &ret,
		}
		reqs = append(reqs, req)  // 将每个 BatchElem 添加到 BatchElem 数组
		rets = append(rets, &ret) // 每个请求的结果引用添加到结果数组中
	}
	err := r.client.GetRpc().BatchCall(reqs) // 传入 BatchElem 数组，发起批量请求
	return rets, err
}

// 根据区块号，获取区块信息
func (r *ETHRPCRequester) GetBlockInfoByNumberService(blockNumber *big.Int) (*model.FullBlock, error) {
	number := fmt.Sprintf("%#x", blockNumber)
	methodName := "eth_getBlockByNumber"
	fullBlock := model.FullBlock{}
	// eth_getBlockByNumber 的第二个参数：
	// 如果是 true，则返回完整的区块信息，false 则 transaction 部分只返回交易hash数组
	err := r.client.client.Call(&fullBlock, methodName, number, true)
	if err != nil {
		return nil, fmt.Errorf("get block info failed! %s", err.Error())
	}
	if fullBlock.Number == "" {
		return nil, fmt.Errorf("block info is empty %s", blockNumber.String())
	}
	return &fullBlock, nil
}

// 获取以太坊最新生成区块的区块号
func (r *ETHRPCRequester) GetLatestBlockNumberService() (*big.Int, error) {
	methodName := "eth_blockNumber"
	number := ""
	// eth_blockNumber 不需要参数
	err := r.client.client.Call(&number, methodName)
	if err != nil {
		return nil, fmt.Errorf("获取最新区块号失败! %s", err.Error())
	}
	ten, _ := new(big.Int).SetString(number[2:], 16)
	return ten, nil
}

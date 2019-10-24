package eth

import (
	"errors"
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

// 单条查询：根据以太坊地址，查询以太坊 eth 的余额
func (r *ETHRPCRequester) GetETHBalanceV(address string, version string) (string, error) {
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
	/*bytes := []byte(result)
	s := hex.EncodeToString(bytes)*/
	ten, _ := new(big.Int).SetString(result[2:], 16)

	return ten.String(), nil

}

//批量查询：根据以太坊地址数组，查询以太坊 eth的余额
func (r *ETHRPCRequester) GetBalances(address []string, version string) ([]string, error) {
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
func (n *NonceManager) GetNonce(address string) *big.Int {
	if n.nonceMemCache == nil {
		n.nonceMemCache = map[string]*big.Int{}
	}
	n.lock.Lock()         // 加锁
	defer n.lock.Unlock() // 当该函数执行完毕，进行解锁
	return n.nonceMemCache[address]
}

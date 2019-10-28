package eth

import (
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
)

type EthRpcClient struct {
	NodeUrl string      //代表节点的url连接
	client  *rpc.Client //代表Rpc客户端实例
}

//传入服务器节点
func NewETHRPCClient(nodeUrl string) *EthRpcClient {
	client := &EthRpcClient{
		NodeUrl: nodeUrl,
	}
	client.initRpc() // 进行初始化rpc客户端句柄实体
	return client
}

// 初始化 rpc 请求实例
func (erc *EthRpcClient) initRpc() {
	// 使用 go-ethereum 库中的 rpc 库来初始化
	// DialHTTP 的意思是使用 http 版本的 rpc 实现方式
	rpcClient, err := rpc.DialHTTP(erc.NodeUrl)
	if err != nil {
		// 初始化失败，终结程序，并将错误信息显示到控制台上面
		errInfo := fmt.Errorf("初始化 rpc client 失败%s", err.Error()).Error()
		panic(errInfo)
	}
	// 初始化成功，将新实例化的 rpc 句柄赋值给我们 ETHRPCClient 结构体里面的
	erc.client = rpcClient
}

//GetRpc方法是为了方便外部能够获取client *rpc.Client 来进行访问
func (erc *EthRpcClient) GetRpc() *rpc.Client {
	if erc.client == nil {
		erc.initRpc()
	}
	return erc.client
}

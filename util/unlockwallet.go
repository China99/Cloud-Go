package util

//解锁钱包
import (
	"errors"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"strings"
)

//全局的保存了已经解锁成功了的钱包map集合变量
var ETHUnlockMap map[string]accounts.Account

//全局对应的keystore实例
var UnlockKs *keystore.KeyStore

//解锁以太坊钱包，传入钱包地址和对应的keystore密码
func UnlockETHWallet(keysDir, address, password string) error {
	if UnlockKs == nil {
		UnlockKs = keystore.NewKeyStore(
			//服务端存储keystore文件的目录
			keysDir,
			keystore.StandardScryptN,
			keystore.StandardScryptP)
		if UnlockKs == nil {
			return errors.New("ks is nil")
		}
	}
	unlock := accounts.Account{Address: common.HexToAddress(address)}
	//ks.Unlock调用keystore.go的解锁函数，解锁出的私钥存储在里面的变量中
	if err := UnlockKs.Unlock(unlock, password); nil != err {
		return errors.New("unlock err:" + err.Error())
	}
	if ETHUnlockMap == nil {
		ETHUnlockMap = map[string]accounts.Account{}
	}
	ETHUnlockMap[address] = unlock //解锁成功，存储
	return nil

}

//签名交易数据结构体 types.Transaction
func SignETHTransaction(address string, transaction *types.Transaction) (*types.Transaction, error) {
	if UnlockKs == nil {
		return nil, errors.New("你需要初始化keystore")
	}
	//
	//keystore存在安装目录的相对路径，进行加密
	account := ETHUnlockMap[address]
	if !common.IsHexAddress(account.Address.String()) {
		//判读当前 地址钱包是否解锁了
		return nil, errors.New("未解锁keystore")
	}
	return UnlockKs.SignTx(account, transaction, nil) //调用签名函数

}

//根据代币的decimal乘以10^decimal后的值
func GetRealDecimalValue(value string, decimal int) string {
	if strings.Contains(value, ".") {
		arr := strings.Split(value, ".")
		if len(arr) != 2 {
			return ""
		}
		//小数位  4.11
		//[11]==2
		num := len(arr[1])
		//18
		left := decimal - num
		return arr[0] + arr[1] + strings.Repeat("0", left)
	} else {
		//整数
		return value + strings.Repeat("0", 18)
	}
}

//小数点左移 decimal位
/*func GetValETH(value []string,decimal int)[]string{
	//strlen := len(value)
	//if strlen>decimal{
	//	for i:=0;i<=strlen-decimal ;i++  {
	//		if i==strlen-decimal{
	//			fa
	//		}
	//	}
	//}
	if len(value)==1{
		str1:=value[0][:decimal]
		str2:=value[0][decimal:]
		str3:=str1+"."+str2
		fmt.Println(str3)
	}
	return
}*/

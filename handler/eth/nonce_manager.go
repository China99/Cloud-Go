package eth

import (
	"math/big"
	"sync"
)

//管理器结构体
type NonceManager struct {
	lock sync.Mutex
	//采用整形大数来存储nonce 区块高度
	nonceMemCache map[string]*big.Int
}

func NewNonceManager() *NonceManager {
	return &NonceManager{
		lock: sync.Mutex{}, //实例化互斥锁
	}
}

//设置nonce
func (n *NonceManager) SetNonce(address string, nonce *big.Int) {
	if n.nonceMemCache == nil {
		n.nonceMemCache = map[string]*big.Int{}
	}
	n.lock.Lock()
	defer n.lock.Unlock()
	n.nonceMemCache[address] = nonce
}

// nonce 进行自增 1
func (n *NonceManager) PlusNonce(address string) {
	if n.nonceMemCache == nil {
		n.nonceMemCache = map[string]*big.Int{}
	}
	n.lock.Lock()         // 加锁
	defer n.lock.Unlock() // 当该函数执行完毕，进行解锁
	oldNonce := n.nonceMemCache[address]
	if oldNonce == nil {
		return
	}
	newNonce := oldNonce.Add(oldNonce, big.NewInt(int64(1)))
	n.nonceMemCache[address] = newNonce
}

package model

//根据文档定义区块信息的结构体
type FullBlock struct {
	Number           string        `json:"number"`     //区块号
	Hash             string        `json:"hash"`       //区块hash
	ParentHash       string        `json:"parentHash"` //前区块hash
	Nonce            string        `json:"nonce"`      //高度
	Sha3Uncles       string        `json:"sha3Uncles"`
	LogsBloom        string        `json:"logsBloom"`
	TransactionsRoot string        `json:"transactionsRoot"`
	ReceiptsRoot     string        `json:"stateRoot"`
	Miner            string        `json:"miner"`
	Difficulty       string        `json:"difficulty"`
	TotalDifficulty  string        `json:"totalDifficulty"`
	ExtraData        string        `json:"extraData"`
	Size             string        `json:"size"`
	GasLimit         string        `json:"gasLimit"`
	GasUsed          string        `json:"gasUsed"`
	Timestamp        string        `json:"timestamp"`
	Uncles           []string      `json:"uncles"`
	Transactions     []Transaction `json:"transactions"`
}

package models

import "time"

type Block struct {
	Id         int
	Number     string
	TimeStamp  time.Time
	TxCount    int
	Hash       string
	ParentHash string
	Producer   string
	Size       string
	GasUsed    uint64
	GasLimit   uint64
	BlockReard string
	ExtraData  string
	Witnesses  string
	Signature  string
	CmtMsges   string
}

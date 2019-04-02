package controllers

import (
	"github.com/vntchain/vnt-explorer/models"
)

type NetController struct {
	BaseController
}

type NetStats struct {
	Height 			uint64
	CurrTps			float32
	TopTps			float32
	TxCount			int64
	AccountCount	int64
	SuperNode		int64
	CandiNode		int64
}

type SearchBody struct {
	Block	*models.Block
	Tx		*models.Transaction
	Account	*models.Account
}

func (this *NetController) Stats() {
	block := &models.Block{}
	block, err := block.Last()
	if err != nil {
		this.ReturnErrorMsg("Failed to get block height: %s", err.Error(), "")
	}

	height := block.Number + 1
	currTps := block.Tps
	topTpsBlock, err := block.TopTpsBlock()
	if err != nil {
		this.ReturnErrorMsg("Failed to get top tps: %s", err.Error(), "")
	}
	topTps := topTpsBlock.Tps

	tx := &models.Transaction{}
	txCount, err := tx.Count("", "", -1, -1, -1)
	if err != nil {
		this.ReturnErrorMsg("Failed to get transaction count: %s", err.Error(), "")
	}

	acct := &models.Account{}
	acctCount, err := acct.Count(-1, -1)

	node := &models.Node{}
	superNode,err := node.Count(1)
	if err != nil {
		this.ReturnErrorMsg("Failed to get super node count: %s", err.Error(), "")
	}

	candiNode,err := node.Count(0)
	if err != nil {
		this.ReturnErrorMsg("Failed to get candidate node count: %s", err.Error(), "")
	}

	stats := &NetStats{
		height,
		currTps,
		topTps,
		txCount,
		acctCount,
		superNode,
		candiNode,
	}

	this.ReturnData(stats, nil)
}
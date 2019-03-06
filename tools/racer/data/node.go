package data

import (
	"github.com/astaxie/beego"
	"github.com/vntchain/vnt-explorer/common"
	"github.com/vntchain/vnt-explorer/common/utils"
	"github.com/vntchain/vnt-explorer/models"

	"math/big"
)

func GetNodes() []*models.Node {
	rpc := common.NewRpc()
	rpc.Method = common.Rpc_GetAllCandidates

	err, resp := utils.CallRpc(rpc)
	if err != nil {
		panic(err.Error())
	}

	nodeList := resp.Result.([]interface{})

	beego.Info("Response body", resp)
	var result []*models.Node
	for _, n := range nodeList {
		node := n.(map[string]interface{})
		address := node["owner"].(string)
		name := node["name"].(string)
		active := node["active"].(bool)
		url := node["url"].(string)
		votes, err := utils.DecodeBig(node["voteCount"].(string))
		if err != nil {
			beego.Error("Get node voteCount err: ", err)
			votes = big.NewInt(0)
		}
		website := node["website"].(string)
		status := 0
		if active {
			status = 1
		}
		nodeValue := models.Node{
			Address: address,
			Vname:   name,
			Home:    website,
			Ip:      url,
			Status:  status,
			Votes:   votes.String(),
		}
		result = append(result, &nodeValue)
	}
	return result
}

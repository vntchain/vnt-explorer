package data

import (
	"math/big"
	"strings"

	"github.com/astaxie/beego"
	"github.com/vntchain/vnt-explorer/common"
	"github.com/vntchain/vnt-explorer/common/utils"
	"github.com/vntchain/vnt-explorer/models"
)

func GetNodes() []*models.Node {
	rpc := common.NewRpc()
	rpc.Method = common.Rpc_GetAllCandidates

	err, resp, _ := utils.CallRpc(rpc)
	if err != nil {
		if err.Error() != "Rpc returned with error: code: -32000, error: empty witness candidates list" {
			beego.Error("Get Node error: ", err)
		}
		return nil
	}

	nodeList := resp.Result.([]interface{})
	totalVotes := big.NewInt(0);

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
		totalBounty, err := utils.DecodeBig(node["totalBounty"].(string))
		if err != nil {
			beego.Error("Get node totalBounty err: ", err)
			votes = big.NewInt(0)
		}
		extractedBounty, err := utils.DecodeBig(node["extractedBounty"].(string))
		if err != nil {
			beego.Error("Get node extractedBounty err: ", err)
			extractedBounty = big.NewInt(0)
		}
		lastExtractTime, err := utils.DecodeBig(node["lastExtractTime"].(string))
		if err != nil {
			beego.Error("Get node lastExtractTime err: ", err)
			lastExtractTime = big.NewInt(0)
		}
		website := node["website"].(string)
		status := 0
		if active {
			status = 1
		}

		// get ip from node url
		tmp := strings.Split(url, "/")
		var ip string
		if len(tmp) > 3 {
			ip = tmp[2]
		}

		totalVotes = totalVotes.Add(totalVotes, votes)

		nodeValue := models.Node{
			Address:         strings.ToLower(address),
			Vname:           name,
			Home:            website,
			Ip:              ip,
			Status:          status,
			Votes:           votes.String(),
			VotesFloat:		 float64(votes.Uint64()),
			TotalBounty:     totalBounty.String(),
			ExtractedBounty: extractedBounty.String(),
			LastExtractTime: lastExtractTime.String(),
			IsAlive: 1,
		}
		result = append(result, &nodeValue)
	}

	votesFloat := float64(totalVotes.Uint64())

	if votesFloat > 0 {
		for _, node := range result {
			node.VotesPercent = float32(node.VotesFloat/votesFloat) * 100
		}
	}
	return result
}

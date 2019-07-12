package main

import (
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/vnt-explorer/models"
	"github.com/vntchain/go-vnt/crypto"
	"fmt"
)

func genNode(addr common.Address, name string) *models.Node {
	return &models.Node{
		Address: addr.String(),
		Vname: name,
		Home: "www.baidu.com",
		Logo: "",
		Ip: "10.0.0.2",
		IsSuper: 0,
		IsAlive: 1,
		Status: 1,
		Votes: "10",
		VotesPercent: 2,
		TotalBounty: "0",
		ExtractedBounty: "0",
		LastExtractTime: "0",
		Longitude: 0,
		Latitude: 0,
		City: "Hangzhou",
		NodeUrl: "/ip4/47.90.248.44/tcp/3001/ipfs/1kHVEheUZE9hxGhiWTsJT2Ft9Wh1R3m18w2UdKjjtDfuKQc",
	}
}

func main() {
	total := 81
	i := 0
	for i < total {
		priv, _ := crypto.GenerateKey()
		pub := crypto.CompressPubkey(&priv.PublicKey)
		address := common.BytesToAddress(pub[:20])
		addressStr := address.String()
		name := fmt.Sprintf("node-%s%s%s%s", addressStr[0], addressStr[8], addressStr[12], addressStr[16])
		node := genNode(address, name)
		if err := node.Insert(); err != nil {
			fmt.Println("Failed to insert node: err", err)
			panic(err)
		}
		i += 1
	}

}

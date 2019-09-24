package data

import (
	"github.com/astaxie/beego"
	"github.com/vntchain/vnt-explorer/common"
	"github.com/vntchain/vnt-explorer/common/utils"
	"github.com/vntchain/vnt-explorer/models"

	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"path"
	"strings"
)

type BpInfo struct {
	Candidate_Name    string
	Candidate_Address string
	Location          Location
	Branding          Branding
}
type Location struct {
	Name      string
	Country   string
	Latitude  float64
	Longitude float64
}

type Branding struct {
	Logo_256  string
	Logo_1024 string
	Logo_Svg  string
}

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

	var result []*models.Node

	if resp.Result != nil {
		nodeList := resp.Result.([]interface{})
		totalVotes := big.NewInt(0)

		beego.Debug("Response body", resp)

		for _, n := range nodeList {
			node := n.(map[string]interface{})
			address := node["owner"].(string)
			name := node["name"].(string)
			registered := node["registered"].(bool)
			bind := node["bind"].(bool)
			url := node["url"].(string)
			votes, err := utils.DecodeBig(node["voteCount"].(string))
			if err != nil {
				beego.Error("Get node voteCount err: ", err)
				votes = big.NewInt(0)
			}
			website := node["website"].(string)
			status := 0
			if registered && bind {
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
				Address:    strings.ToLower(address),
				Vname:      name,
				Home:       website,
				Ip:         ip,
				NodeUrl:    url,
				Status:     status,
				Votes:      votes.String(),
				VotesFloat: float64(votes.Uint64()),
				Latitude:   360,
				Longitude:  360,
			}
			result = append(result, &nodeValue)
		}

		votesFloat := float64(totalVotes.Uint64())

		if votesFloat > 0 {
			for _, node := range result {
				node.VotesPercent = float32(node.VotesFloat/votesFloat) * 100
			}
		}
	}

	return result
}

func GetBpInfo(website string) (bp *BpInfo) {
	body, err := utils.CallApi(website, nil)
	if err != nil {
		beego.Error("Faile to CallApi ", website, " error  ", err.Error())
		return nil
	}

	bp = &BpInfo{}
	err = json.Unmarshal(body, bp)
	if err != nil {
		beego.Error("Failed to unmarshal bpInfo: %s", err.Error())
		return nil
	}
	if bp.Location.Longitude < -180 || bp.Location.Longitude > 180 ||
		bp.Location.Latitude < -90 || bp.Location.Latitude > 90 {
		bp.Location.Longitude = 360
		bp.Location.Latitude = 360
	}
	return
}

func GetLogo(imgUrl, address string) {
	imgName := path.Base(imgUrl)
	imgDir := path.Join(common.IMAGE_PATH, address)
	imgPath := path.Join(imgDir, imgName)
	if exists, _, _ := FileExists(imgPath); exists {
		return
	}
	if exists, _, _ := FileExists(imgDir); !exists {
		err := os.MkdirAll(imgDir, 0711)
		if err != nil {
			beego.Error("Failed to create image dir: %s", imgDir)
		}
	}
	res, err := http.Get(imgUrl)
	if err != nil {
		beego.Error("Failed to download logo of %s: %s", address, err.Error())
		return
	}
	if res == nil || res.Body == nil {
		return
	}

	defer res.Body.Close()
	file, err := os.Create(imgPath)
	if err != nil {
		beego.Error("Failed to create logo file of %s: %s", address, err.Error())
		return
	}

	_, err = io.Copy(file, res.Body)
	if err != nil {
		beego.Error("Failed to write logo of %s: %s", address, err.Error())
	}
}

func FileExists(filePath string) (bool, int64, error) {
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false, 0, nil
	}
	if err != nil {
		msg := fmt.Sprintf("error [%s] checking if file [%s] exists", err, filePath)
		return false, 0, errors.New(msg)
	}
	return true, fileInfo.Size(), nil
}

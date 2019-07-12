package main

import (
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/vntchain/vnt-explorer/models"
	"github.com/vntchain/vnt-explorer/tools/racer/data"
	"github.com/vntchain/vnt-explorer/tools/racer/token"
	"runtime/debug"
	"sync"
	"github.com/vntchain/vnt-explorer/common"
)

func main() {
	common.InitLogLevel()
	registerElectionContract()

	for {
		doSync()
	}
}

func doSync() {
	defer func() {
		if r := recover(); r != nil {
			beego.Error("Error happened:", r)
			debug.PrintStack()
			time.Sleep(2 * time.Second)
		}
	}()

	var rmtHgt, localHgt int64
	rmtHgt, localHgt = checkHeight()

	// localHgt = 14
	//rmtHgt = 25913
	beego.Info(fmt.Sprintf("Local height: %d, rmtHeight: %d", localHgt, rmtHgt))
	if localHgt >= rmtHgt {
		beego.Info("No more blocks.")
		time.Sleep(2 * time.Second)
		return
	}

	// Set the block sync batch
	if rmtHgt-localHgt > data.BatchSize {
		rmtHgt = localHgt + data.BatchSize
	}

	for localHgt < rmtHgt {
		localHgt = localHgt + 1
		beego.Info(fmt.Sprintf("Will sync block %d", localHgt))
		data.PostBlockTask(data.NewBlockTask(localHgt))
	}

	data.PostNodesTask(data.NewNodesTask())

	Wait()
}

func Wait() {
	var wait = 3
	for {
		blockQueued, blockActive := data.BlockPool.Routines()
		blockInsertPoolQ, blockInsertPoolA := data.BlockInsertPool.Routines()
		txQueued, txActive := data.TxPool.Routines()
		accoutExtQueued, accountExtActive := data.AccountExtractPool.Routines()
		accoutQueued, accountActive := data.AccountPool.Routines()
		witQueued, witActive := data.WitnessesPool.Routines()
		nodeQueded, nodeActive := data.NodePool.Routines()
		nodeInfoQueded, nodeInfoActive := data.NodeInfoPool.Routines()
		logoQueded, logoActive := data.NodeInfoPool.Routines()
		//beego.Info("Work Pool status:")
		//beego.Info("Block Pool:", "Queued=", blockQueued, ",Active=", blockActive)
		//beego.Info("Block Insert Pool:", "Queued=", blockInsertPoolQ, ",Active=", blockInsertPoolA)
		//beego.Info("Tx Pool:", "Queued=", txQueued, ",Active=", txActive)
		//beego.Info("Account Extract Pool:", "Queued=", accoutExtQueued, ",Active=", accountExtActive)
		//beego.Info("Account Pool:", "Queued=", accoutQueued, ",Active=", accountActive)
		//beego.Info("Witness Pool:", "Queued=", witQueued, ",Active=", witActive)
		//beego.Info("Node Pool:", "Queued=", nodeQueded, ",Active=", nodeActive)

		if blockQueued+blockActive == 0 &&
			blockInsertPoolQ+blockInsertPoolA == 0 &&
			txQueued+txActive == 0 &&
			accoutExtQueued+accountExtActive == 0 &&
			accoutQueued+accountActive == 0 &&
			witQueued+witActive == 0 &&
			nodeQueded+nodeActive == 0 &&
			nodeInfoQueded+nodeInfoActive == 0 &&
			logoQueded+logoActive == 0 {
			if wait == 0 {
				data.AccountMap = sync.Map{}
				token.TokenMap = sync.Map{}
				break
			} else {
				wait--
				time.Sleep(500 * time.Millisecond)
				continue
			}

			//time.Sleep(1 * time.Second)
		}

		wait = 3
		time.Sleep(500 * time.Millisecond)
	}
}

func checkHeight() (int64, int64) {
	rmtHgt := data.GetRemoteHeight()
	localHgt, _ := data.GetLocalHeight()

	if localHgt > rmtHgt {
		msg := fmt.Sprintf("Local height %d is bigger than remote height: %d, please check your remote node", localHgt, rmtHgt)
		beego.Error(msg)
		panic(msg)
	}

	return rmtHgt, localHgt
}

func registerElectionContract() {
	electionAddr := "0x0000000000000000000000000000000000000009"
	election := &models.Account{}
	election, err := election.Get(electionAddr)
	if err != nil {
		election = &models.Account{
			Address:        electionAddr,
			Balance:        "0",
			TxCount:        0,
			IsContract:     true,
			ContractName:   "election",
			TokenAcctCount: "0",
			TokenAmount:    "0",
		}
		if err = election.Insert(); err != nil {
			msg := fmt.Sprintf("Failed to insert election contract account: %s", err.Error())
			panic(msg)
		}
	} else if election.ContractName == "" {
		election.ContractName = "election"
		if err = election.Update(); err != nil {
			msg := fmt.Sprintf("Failed to update election contract account: %s", err.Error())
			panic(msg)
		}
	}
}

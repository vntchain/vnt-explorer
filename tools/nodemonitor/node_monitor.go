package main

import (
	"context"
	"fmt"
	"time"

	"encoding/binary"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-host"
	"github.com/libp2p/go-libp2p-net"
	"github.com/libp2p/go-libp2p-peerstore"
	"github.com/multiformats/go-multiaddr"
	"github.com/vntchain/go-vnt/rlp"
	"github.com/vntchain/vnt-explorer/common"
	"github.com/vntchain/vnt-explorer/models"
	"io"
	"math/big"
)

const (
	MessageHeaderLength = 5
	ProtocolID          = "vnt"
	StatusMsg           = 0x00
	ProtocolMaxMsgSize  = 10 * 1024 * 1024 // Maximum cap on the size of a protocol message
	HashLength          = 32
)

const (
	ErrMsgTooLarge = iota
	ErrDecode
	ErrInvalidMsgCode
	ErrProtocolVersionMismatch
	ErrNetworkIdMismatch
	ErrGenesisBlockMismatch
	ErrNoStatusMsg
	ErrExtraStatusMsg
	ErrSuspendedPeer
	ErrLowTD
)

type activeStatus struct {
	nodeUrl string
	active  bool
}

var (
	genesisHash string
	nodePool    chan string
	resPool     chan activeStatus
	td          uint64
	nodeMap     map[string]*models.Node
)
var interval, intervalErr = beego.AppConfig.Int("supernode::interval")
var sourcePort = beego.AppConfig.String("supernode::p2p_port")

type MessageType uint64
type errCode int
type Hash [HashLength]byte

// Msg message struct
type Msg struct {
	Header MsgHeader
	Body   MsgBody
}

// MsgHeader store the size of MsgBody
type MsgHeader [MessageHeaderLength]byte

// MsgBody message body
type MsgBody struct {
	ProtocolID  string //Protocol name
	Type        MessageType
	ReceivedAt  time.Time
	PayloadSize uint32
	Payload     io.Reader
}

// GetBodySize get message body size in uint32
func (msg *Msg) GetBodySize() uint32 {
	header := msg.Header
	bodySize := binary.LittleEndian.Uint32(header[:])
	return bodySize
}

// Decode using json unmarshal decode msg payload
func (msg Msg) Decode(val interface{}) error {
	s := rlp.NewStream(msg.Body.Payload, uint64(msg.Body.PayloadSize))
	err := s.Decode(val)
	if err != nil {
		beego.Error("Decode()", "err", err, "message type", msg.Body.Type, "payload size", msg.Body.PayloadSize)
		return err
	}
	return nil
}

type statusData struct {
	ProtocolVersion uint32
	NetworkId       uint64
	TD              *big.Int
	CurrentBlock    Hash
	GenesisBlock    Hash
}

func errResp(code errCode, format string, v ...interface{}) error {
	return fmt.Errorf("%v - %v", code, fmt.Sprintf(format, v...))
}

func readData(s net.Stream, genesis string, localTD *big.Int) error {
	defer s.Reset()
	for {
		// 读取消息
		msgHeaderByte := make([]byte, MessageHeaderLength)
		_, err := io.ReadFull(s, msgHeaderByte)
		if err != nil {
			return fmt.Errorf("HandleStream read msg header error : %s", err)
		}
		bodySize := binary.LittleEndian.Uint32(msgHeaderByte)

		msgBodyByte := make([]byte, bodySize)
		_, err = io.ReadFull(s, msgBodyByte)
		if err != nil {
			return fmt.Errorf("HandleStream read msg Body error: %s", err)
		}
		msgBody := &MsgBody{Payload: &rlp.EncReader{}}
		err = json.Unmarshal(msgBodyByte, msgBody)
		if err != nil {
			return fmt.Errorf("HandleStream unmarshal msg Body error: %s", err)
		}
		msgBody.ReceivedAt = time.Now()

		// 传递给msger
		var msgHeader MsgHeader
		copy(msgHeader[:], msgHeaderByte)

		msg := Msg{
			Header: msgHeader,
			Body:   *msgBody,
		}
		if msgBody.ProtocolID == ProtocolID {
			return readStatus(msg, genesis, localTD)
		}
	}
	return nil
}

func readStatus(msg Msg, genesis string, localTD *big.Int) (err error) {
	var status *statusData
	if msg.Body.Type != StatusMsg {
		return errResp(ErrNoStatusMsg, "first msg has code %x (!= %x)", msg.Body.Type, StatusMsg)
	}
	size := msg.GetBodySize()
	if size > ProtocolMaxMsgSize {
		return errResp(ErrMsgTooLarge, "%v > %v", size, ProtocolMaxMsgSize)
	}
	// Decode the handshake and make sure everything matches
	if err := msg.Decode(&status); err != nil {
		return errResp(ErrDecode, "msg %v: %v", msg, err)
	}

	if genesis != "" && fmt.Sprintf("0x%x", status.GenesisBlock) != genesis {
		return errResp(ErrGenesisBlockMismatch, "%x (!= %x)", status.GenesisBlock[:8], genesis[:8])
	}
	if status.TD.Cmp(localTD) < 0 {
		return errResp(ErrLowTD, "remote's TD %v (< local TD %v)", status.TD, localTD)
	}
	return nil
}

func pingNode(host host.Host, nodeUrl string) { // Turn the destination into a multiaddr.
	maddr, err := multiaddr.NewMultiaddr(nodeUrl)
	if err != nil {
		beego.Error("ping Node, get maddr error ", err, " nodeUrl ", nodeUrl)
		resPool <- activeStatus{nodeUrl, false}
		return
	}

	// Extract the peer ID from the multiaddr.
	info, err := peerstore.InfoFromP2pAddr(maddr)
	if err != nil {
		beego.Error("ping Node, InfoFromP2pAddr error ", err, " maddr ", err)
		resPool <- activeStatus{nodeUrl, false}
		return
	}

	// Add the destination's peer multiaddress in the peerstore.
	// This will be used during connection and stream creation by libp2p.
	host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	// Start a stream with the destination.
	// Multiaddress of the destination peer is fetched from the peerstore using 'peerId'.
	s, err := host.NewStream(context.Background(), info.ID, "/p2p/1.0.0")
	if err != nil {
		beego.Error("NewStream error: ", err, " nodeUrl ", nodeUrl)
		resPool <- activeStatus{nodeUrl, false}
		return
	} else if err := readData(s, genesisHash, big.NewInt(0).SetUint64(td)); err != nil {
		beego.Error("read data error ", err, " nodeUrl ", nodeUrl)
		resPool <- activeStatus{nodeUrl, false}
		return
	}
	resPool <- activeStatus{nodeUrl, true}
}

func pingManager(host host.Host) {
	for {
		select {
		case url := <-nodePool:
			go pingNode(host, url)
		}
	}
}

func getAllNodes() {
	if genesisHash == "" {
		getGenesis()
	}
	getTd()

	n := &models.Node{}
	nodes, err := n.All()
	if err != nil {
		beego.Error("Get Nodes from db error ", err)
		return
	}

	for _, node := range nodes {
		nodeMap[node.NodeUrl] = node
		nodePool <- node.NodeUrl
	}
}

func updateDB() {
	for {
		select {
		case res := <-resPool:
			if node, exists := nodeMap[res.nodeUrl]; exists && node != nil && res.active != (node.IsAlive == 1) {
				node.IsAlive = 1 - node.IsAlive
				nodeMap[res.nodeUrl].Insert()
			}
		}
	}
}

func main() {
	if intervalErr != nil {
		interval = common.DefaultNodeInterval
	}

	// 0.0.0.0 will listen on any interface device.
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", sourcePort))

	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	host, err := libp2p.New(
		context.Background(),
		libp2p.ListenAddrs(sourceMultiAddr),
	)
	if err != nil {
		panic(err)
	}

	beego.Info("This node's multiaddresses:")
	for _, la := range host.Addrs() {
		beego.Info("-", la)
	}

	nodePool = make(chan string, 64)
	resPool = make(chan activeStatus, 64)
	nodeMap = make(map[string]*models.Node)

	go pingManager(host)
	go getAllNodes()
	go updateDB()
	t := time.Tick(time.Second * time.Duration(interval))
	for range t {
		go getAllNodes()
	}
}

func getGenesis() {
	genesis := &models.Block{}
	if genesis, err := genesis.GetByNumber(0); err == nil && genesis != nil {
		genesisHash = genesis.Hash
	}
}

func getTd() {
	block := &models.Block{}
	if last, err := block.Last(); err == nil && last != nil {
		td = last.Number + 1
	}
}

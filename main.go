package main

import (
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	log "github.com/sirupsen/logrus"
)

type Snapshot struct {
	Timestamp    int                    `json:"timestamp"`
	TotalNodes   int                    `json:"total_nodes"`
	LatestHeight int32                  `json:"latest_height"`
	Nodes        map[string]interface{} `json:"nodes"`
}

const PVER = wire.ProtocolVersion
const BTCNET = wire.MainNet
const THREADS = 500
const TIMEOUT = 30

var latest_height int32
var header = wire.BlockHeader{
	Version:    int32(rand.Int()),
	PrevBlock:  chainhash.Hash{},
	MerkleRoot: chainhash.Hash{},
	Timestamp:  time.Now(),
	Bits:       uint32(rand.Int()),
	Nonce:      uint32(rand.Int()),
}

// Fetches the latest nodes from BitNodes, and decodes into the Snapshot struct
func getNodes(snapshot *Snapshot) error {
	// Fetch latest bitcoin nodes.
	log.Info("Fetching latest node snapshot")
	response, err := http.Get("https://bitnodes.io/api/v1/snapshots/latest/")
	if err != nil {
		return err
	}

	// Parse response
	if err := json.NewDecoder(response.Body).Decode(&snapshot); err != nil {
		return err
	}
	return nil
}

func nodeHandshake(conn net.Conn) error {
	// Create Addr objects for VERSION handshake
	you := wire.NewNetAddress(conn.RemoteAddr().(*net.TCPAddr), wire.SFNodeNetwork)
	me := wire.NewNetAddress(conn.LocalAddr().(*net.TCPAddr), wire.SFNodeNetwork)

	// Initiate handshake with node
	nonce := rand.Uint64()
	version := wire.NewMsgVersion(me, you, nonce, latest_height)
	err := wire.WriteMessage(conn, version, PVER, BTCNET)
	if err != nil {
		return err
	}

	time.Sleep(time.Second)

	// Send verack
	verack := wire.NewMsgVerAck()
	err = wire.WriteMessage(conn, verack, PVER, BTCNET)
	if err != nil {
		return err
	}
	return nil
}

func waitUntilConnectionCloses(conn net.Conn, deadline time.Time) error {
	// End function when connection is closed or 30 sec timeout
	for {
		_, err := conn.Read(make([]byte, 1))
		if err != nil {
			if err == io.EOF || err == net.ErrClosed {
				return nil
			}
			if errors.Is(err, os.ErrDeadlineExceeded) {
				return errors.New("Read buffer timeout")
			}
		}
	}
}

func discourageIP(ip string) {
	log.Debug("Starting: " + ip)
	deadline := time.Now().Add(TIMEOUT * time.Second)
	// Validate IP
	i := ip[:strings.LastIndex(ip, ":")]
	if_ip := net.ParseIP(i)
	if if_ip == nil {
		log.Debug("Invalid IP: " + ip)
		return
	}
	// Open a connection to the Bitcoin node.
	conn, err := net.DialTimeout("tcp", ip, time.Second*10)
	if err != nil {
		log.Debug("TCP Connection failure: " + err.Error())
		return
	}
	defer conn.Close()

	conn.SetDeadline(deadline)

	// Node handshake is required before the node will accept any other messages from peer
	err = nodeHandshake(conn)
	if err != nil {
		log.Debug("Node Handshake failure: " + err.Error())
		return
	}

	// Send phony header
	headers := wire.NewMsgHeaders()
	headers.AddBlockHeader(&header)
	err = wire.WriteMessage(conn, headers, PVER, BTCNET)
	if err != nil {
		log.Debug("Payload failure: " + err.Error())
		return
	}

	// Wait for the remote server to close the connection.
	err = waitUntilConnectionCloses(conn, deadline)
	if err != nil {
		log.Debug("TCP Close failure: " + err.Error())
		return
	}
	log.Info("Successfully discouraged: " + ip)
}

func main() {
	// Setup logging
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.InfoLevel)

	// Fetch nodes
	var snapshot Snapshot
	err := getNodes(&snapshot)
	if err != nil {
		log.Fatal("Fetch nodes error: " + err.Error())
	}

	// Initialise some runtime constants for the payload
	latest_height = snapshot.LatestHeight
	hash := []byte("Change The Code")
	header.PrevBlock.SetBytes(hash)
	header.MerkleRoot.SetBytes(hash)

	log.Info("Discouraging nodes...")

	// Iterate over each node for discouragement, using goroutines and WaitGroup
	wg := sync.WaitGroup{}
	wg.Add(len(snapshot.Nodes))
	guard := make(chan struct{}, THREADS)

	for ip := range snapshot.Nodes {

		guard <- struct{}{}

		go func(ip string) {
			defer func() {
				wg.Done()
				<-guard
			}()

			discourageIP(ip)
		}(ip)
	}
	wg.Wait()

	log.Info("Successfully discouraged your public IP from all IPv4 Bitcoin nodes for 24 hours.\nThanks for caring for the environment.\nPlease run again tomorrow")
}

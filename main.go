package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

type Snapshot struct {
	Timestamp    int                    `json:"timestamp"`
	TotalNodes   int                    `json:"total_nodes"`
	LatestHeight int32                  `json:"latest_height"`
	Nodes        map[string]interface{} `json:"nodes"`
}

var pver = wire.ProtocolVersion
var btcnet = wire.TestNet3
var latest_height int32
var header = wire.BlockHeader{
	Version:    int32(rand.Int()),
	PrevBlock:  chainhash.Hash{},
	MerkleRoot: chainhash.Hash{},
	Timestamp:  time.Now(),
	Bits:       0x1d00ffff,
	Nonce:      uint32(rand.Int()),
}

// Fixed function that fetches the latest nodes from BitNodes, and decodes into the Snapshot struct
func getNodes(snapshot *Snapshot) {
	// Fetch latest bitcoin nodes.
	fmt.Println("Fetching latest node snapshot")
	response, err := http.Get("https://bitnodes.io/api/v1/snapshots/latest/")
	if err != nil {
		fmt.Println(err)
	}

	// Parse response
	defer response.Body.Close()
	if err := json.NewDecoder(response.Body).Decode(&snapshot); err != nil {
		fmt.Println(err)
	}
}

func waitUntilConnectionCloses(conn net.Conn) error {
	// End function when connection is closed or 30 sec timeout
	deadline := time.Now().Add(30 * time.Second)
	for {
		_, err := conn.Read(make([]byte, 1))
		if err != nil {
			if err == io.EOF || err == net.ErrClosed {
				fmt.Println("Success: Connection closed by peer")
				return nil
			}
			if time.Now().After(deadline) {
				return errors.New("Connection timed out")
			}
			return err
		}
	}
}

func discourageIP(ip string) {
	// Open a connection to the Bitcoin node.
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// Create Addr objects for VERSION handshake
	you := wire.NewNetAddress(conn.RemoteAddr().(*net.TCPAddr), wire.SFNodeNetwork)
	me := wire.NewNetAddress(conn.LocalAddr().(*net.TCPAddr), wire.SFNodeNetwork)

	// Initiate handshake with node
	nonce := rand.Uint64()
	version := wire.NewMsgVersion(me, you, nonce, latest_height)
	err = wire.WriteMessage(conn, version, pver, btcnet)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Receive version response from peer OPTIONAL?
	// response, _, err := wire.ReadMessage(conn, pver, btcnet)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(response)

	// Send verack
	verack := wire.NewMsgVerAck()
	err = wire.WriteMessage(conn, verack, pver, btcnet)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Send phony header
	headers := wire.NewMsgHeaders()
	headers.AddBlockHeader(&header)
	err = wire.WriteMessage(conn, headers, pver, btcnet)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Wait for the remote server to close the connection.
	err = waitUntilConnectionCloses(conn)
	if err != nil {
		fmt.Println(err)
		return
	}

}

func main() {
	var snapshot Snapshot
	getNodes(&snapshot)

	latest_height = snapshot.LatestHeight

	fmt.Println("Discouraging nodes...")

	// // Iterate over each node for discouragement
	// for ip := range snapshot.Nodes {
	// 	discourageIP(ip)
	// }
}

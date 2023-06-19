package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Snapshot struct {
	Timestamp    int                    `json:"timestamp"`
	TotalNodes   int                    `json:"total_nodes"`
	LatestHeight int                    `json:"latest_height"`
	Nodes        map[string]interface{} `json:"nodes"`
}

// Fixed function that fetches the latest nodes from BitNodes, and decodes into the Snapshot struct
func getNodes(snapshot *Snapshot) {
	// Fetch latest bitcoin nodes.
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

func discourageIP(ip string) {
	fmt.Println(ip)
}

func main() {
	var snapshot Snapshot
	getNodes(&snapshot)

	fmt.Printf("Fetched latest nodes: %b total nodes", snapshot.TotalNodes)

	// Iterate over each node for discouragement
	for ip := range snapshot.Nodes {
		discourageIP(ip)
	}
}

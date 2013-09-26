package main

/* Things to expose:
cjdns uptime
cjdns version

*/

import (
	"fmt"

	// This is the 'refactor' branch of go-cjdns. Just clone the repo to this
	// folder and switch branches.
	"github.com/inhies/go-cjdns-refactor/cjdns"
	"github.com/inhies/go-log"
	"github.com/kylelemons/godebug/pretty"
	"net/http"
	"os"
	"runtime"
	"time"
)

// Peer data.
type Peer struct {
	PublicKey   string
	State       cjdns.PeerState
	IsIncoming  bool
	BytesIn     int64
	BytesOut    int64
	Last        time.Time
	SwitchLabel string

	IPv6       string
	RateIn     float64
	RateOut    float64
	LastUpdate time.Time
}

// Data that has been gathered from cjdns
var Data struct {
	// Data on peers with the public key as the map key so we can easily find
	// a specific peer. This might change, because the JSON output is ugly.
	Peers map[string]Peer

	// Data on this specific node.
	Node struct {
		Memory int64
		Angel  struct {
			Uptime        Duration
			PercentCPU    float64
			PercentMemory float64
		}

		Core struct {
			Uptime        Duration
			PercentCPU    float64
			PercentMemory float64
		}

		RateIn  float64
		RateOut float64

		// The difference of all bytes received and sent
		BytesIn  int64
		BytesOut int64

		angelPID []byte
		corePID  []byte
	}
	CjdnsConn *cjdns.Conn `json:"-"`
}

// Create the programs configuration struct
var SystemConfig = new(configuration)

// For debugging structs
var Pretty = &pretty.Config{PrintStringers: true}

// Program-wide logger
var l *log.Logger

func init() {
	//makeConfig()
	//return

	// read config
	if err := parseConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Convert between the config file log level, which could be an int or a
	// string
	logLevel, err := log.ParseLevel(SystemConfig.Log.RawLogLevel)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Set the log level and start the logger
	SystemConfig.Log.LogLevel = logLevel
	l, err = log.NewLevel(logLevel, SystemConfig.Log.IncLogLevel, os.Stdout,
		"", log.Lshortfile|log.Ldate|log.Lmicroseconds)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Start logging to stdout
	l = log.New(os.Stdout, "", log.Lshortfile|log.Ldate|log.Lmicroseconds)
}

func main() {

	// Connect to cjdns using the default ~/.cjdnsadmin file
	conn, err := cjdns.Connect(nil)
	if err != nil {
		l.Fatalln(err)
	}

	// Set the cjdns connection and create the peer data map
	Data.CjdnsConn = conn
	Data.Peers = make(map[string]Peer)

	// Launch the peer stat fetching routine.
	go peerStatLoop()

	// Start the HTTP JSON API if enabled.
	if SystemConfig.Access.JSONApi.Enabled {
		l.Infoln("Starting HTTP JSON API")
		http.HandleFunc("/peers/", peerStatsHandler)
		http.HandleFunc("/node/", nodeStatsHandler)
		http.HandleFunc("/all/", allStatsHandler)
		http.HandleFunc("/static/", assetsHandler)
		http.ListenAndServe(SystemConfig.Access.JSONApi.Addr, nil)
	}

	l.Debugln("All enabled access methods have been started")
	// Enter in to an infinite loop
	for {
		runtime.Gosched()
	}
}

package main

/* Things to expose:
cjdns uptime
cjdns version

*/

import (
	"fmt"
	"github.com/inhies/go-cjdns/admin"
	"github.com/inhies/go-cjdns/key"
	"github.com/inhies/go-log"
	"github.com/kylelemons/godebug/pretty"
	"net"
	"os"
	"runtime"
	"time"
)

// Peer data.
type Peer struct {
	PublicKey   *key.Public
	State       string
	IsIncoming  bool
	BytesIn     int
	BytesOut    int
	Last        time.Time
	SwitchLabel *admin.Path

	IPv6       net.IP
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
		Memory int
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
		BytesIn  int
		BytesOut int

		angelPID []byte
		corePID  []byte
	}
	CjdnsConn *admin.Conn `json:"-"`
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

	/*
	   // For some reason this results in an infinite loop, with the same error
	   // message displayed. It LOOKS right but I'm not sure why this happens
	   var conn *cjdns.Conn
	   var err error
	   for conn, err = cjdns.Connect(nil); err != nil; {
	       l.Emergln(err)
	       time.Sleep(1 * time.Second)
	   }
	*/

	// This works but I don't like it. I wish the above code would play nice...
connect:
	conn, err := admin.Connect(nil)
	if err != nil {
		l.Emergln(err)
		time.Sleep(1 * time.Second)
		goto connect
	}

	// Set the cjdns connection and create the peer data map
	Data.CjdnsConn = conn
	Data.Peers = make(map[string]Peer)

	// Launch the peer stat fetching routine.
	go peerStatLoop()

	// Start web server if enabled
	Serve()

	l.Debugln("All enabled access methods have been started")
	// Enter in to an infinite loop
	for {
		runtime.Gosched()
	}
}

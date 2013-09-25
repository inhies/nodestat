package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"time"
)

// Updates the cjdns statistics in the global Data struct.
func updateCjdnsStats() (err error) {
	Data.Node.Memory, err = Data.CjdnsConn.Memory()
	if err != nil {
		return
	}

	results, err := exec.Command("pidof", "cjdroute").Output()
	if err != nil {
		return
	}

	rawPids := bytes.Fields(results) //Split(results, []byte("\n"))

	if len(rawPids) > 2 {
		return fmt.Errorf("Too many cjdroute instances running")
	} else if len(rawPids) < 2 {
		return fmt.Errorf("Not enough cjdroute instances running")
	}
	for _, p := range rawPids {

		// Get process uptime, percent cpu, percent memory, and arguments
		results, err = exec.Command("ps", "-p", string(p), "-o",
			"etimes=,pcpu=,pmem=,args=").CombinedOutput()

		if err != nil {
			return
		}

		// We should have at least 6 fields
		data := bytes.Fields(results)
		if len(data) < 6 {
			return fmt.Errorf("Not enough information returned from ps")
		}

		// Parse the uptime in seconds
		tempUptime, err := time.ParseDuration(string(bytes.TrimSpace(data[0])) + "s")

		if err != nil {
			return err
		}

		// The second argument returned from 'args=' should be angel or core
		switch string(data[4]) {
		case "core":
			n := &Data.Node.Core
			n.Uptime = Duration(tempUptime)
			n.PercentCPU, err = strconv.ParseFloat(string(data[1]), 64)
			if err != nil {
				return err
			}
			n.PercentMemory, err = strconv.ParseFloat(string(data[2]), 64)
			if err != nil {
				return err
			}

		case "angel":
			n := &Data.Node.Angel
			n.Uptime = Duration(tempUptime)
			n.PercentCPU, err = strconv.ParseFloat(string(data[1]), 64)
			if err != nil {
				return err
			}
			n.PercentMemory, err = strconv.ParseFloat(string(data[2]), 64)
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("Can not determine if PID belongs to angel or core.")
		}
	}
	return
}

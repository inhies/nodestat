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

	rawPids := bytes.Fields(results)

	if len(rawPids) > 2 {
		return fmt.Errorf("Too many cjdroute instances running")
	} else if len(rawPids) < 2 {
		return fmt.Errorf("Not enough cjdroute instances running")
	}
	for _, p := range rawPids {
		// Get process uptime, percent cpu, percent memory, and arguments.
		// For some reason I have to use multiple -o to get ps to play nice...
		results, err = exec.Command("ps", "-p", string(p), "-o",
			"etime=", "-o", "pcpu=", "-o", "pmem=", "-o",
			"args=").CombinedOutput()

		if err != nil {
			return
		}

		// We should have at least 6 fields
		data := bytes.Fields(results)
		if len(data) < 6 {
			l.Debugln(data)
			return fmt.Errorf("Not enough information returned from ps")
		}

		// Parse the uptime
		rawDur := bytes.TrimSpace(data[0])
		var splitTime [][]byte
		var hours uint64
		if bytes.Contains(rawDur, []byte("-")) {
			// Number of days was included in the response, so we'll split it
			// off and convert to hours
			splitDay := bytes.Split(rawDur, []byte("-"))
			var days uint64
			days, err = strconv.ParseUint(string(splitDay[0]), 10, 64)
			if err != nil {
				return
			}
			hours = days * 24

			splitTime = bytes.Split(splitDay[1], []byte(":"))
		} else {
			// We just have hh:mm:ss data, so split it up.
			splitTime = bytes.Split(rawDur, []byte(":"))
		}

		// Convert the hours to a uint64 so we can add them to any hours we may
		// have from the day conversion
		var rawHours uint64
		rawHours, err = strconv.ParseUint(string(splitTime[0]), 10, 64)
		if err != nil {
			return
		}
		hours += rawHours

		// Manually create the duration string for time.ParseDuration
		var duration string
		switch len(splitTime) {
		case 3:
			duration = strconv.FormatUint(hours, 10) + "h" +
				string(splitTime[1]) + "m" + string(splitTime[2]) + "s"
		case 2:
			duration = string(splitTime[0]) + "m" + string(splitTime[1]) + "s"
		default:
			err = fmt.Errorf("Unable to parse cjdns uptime")
			return
		}

		tempUptime, err := time.ParseDuration(duration)

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

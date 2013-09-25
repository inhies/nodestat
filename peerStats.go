package main

import (
	"github.com/inhies/go-cjdns-refactor/cjdns"
	"time"
)

func peerStatLoop() {
	for {
		// Sleep until its time to update
		time.Sleep(time.Duration(SystemConfig.DataSources.CjdnsPeers.Interval))

		// Get peer stats from cjdns
		results, err := Data.CjdnsConn.InterfaceController_peerStats(0)
		if err != nil {
			l.Noticeln(err)
			return
		}

		Data.Node.RateIn, Data.Node.RateOut = 0, 0
		Data.Node.BytesIn, Data.Node.BytesOut = int64(0), int64(0)
		peerUpdate := make(map[string]Peer)
		// loop through the results and update peer statistics while checking
		// for any peers that are no longer there.
		for _, peer := range results {
			// Create the peer in the map if it's not there yet and calculate
			// the IPv6 address. If it is there, just get the IPv6 address.
			var IPv6 string
			if _, ok := Data.Peers[peer.PublicKey]; !ok {
				IPv6, err = cjdns.PubKeyToIP(peer.PublicKey)
				if err != nil {
					l.Noticeln(err)
					continue
				}
				//IPv6 = net.ParseIP(IP)
				Data.Peers[peer.PublicKey] = Peer{
					IPv6: IPv6,
				}
			} else {
				IPv6 = Data.Peers[peer.PublicKey].IPv6
			}

			// Calculate upload and download rate
			newRateIn := float64(peer.BytesIn-Data.Peers[peer.PublicKey].BytesIn) / (time.Since(Data.Peers[peer.PublicKey].LastUpdate).Seconds())
			newRateOut := float64(peer.BytesOut-Data.Peers[peer.PublicKey].BytesOut) / (time.Since(Data.Peers[peer.PublicKey].LastUpdate).Seconds())

			// Update the peer statistics
			peerUpdate[peer.PublicKey] = Peer{
				PublicKey:   peer.PublicKey,
				IPv6:        Data.Peers[peer.PublicKey].IPv6, // save the same IP
				LastUpdate:  time.Now(),
				Last:        peer.Last,
				SwitchLabel: peer.SwitchLabel,
				IsIncoming:  peer.IsIncoming,
				State:       peer.State,

				BytesIn:  peer.BytesIn,
				RateIn:   newRateIn,
				BytesOut: peer.BytesOut,
				RateOut:  newRateOut,
			}
			Data.Node.RateIn += newRateIn
			Data.Node.RateOut += newRateOut
			Data.Node.BytesIn += peer.BytesIn
			Data.Node.BytesOut += peer.BytesOut

		}

		// Copy the new peer data to the persistant Data.Peers, this should drop
		// any non existant peers.
		Data.Peers = peerUpdate

		/* for _, p := range Data.Peers {
			//fmt.Println(
			//p.IPv6.String(),
			//p.RateIn.String()+"/s", "/", p.RateOut.String()+"/s", "In/Out", "Status:", p.State.String())
		}
		*/
		/*
			output, err := json.MarshalIndent(Data.Peers, "", "\t")
			if err != nil {
				fmt.Println(err)
				//Pretty.Print(Data.Peers)
				return
			}

			fmt.Println(string(output))
		*/
		//fmt.Println("Totals:", totalIn.String()+"/s", "/", totalOut.String()+"/s", "In/Out")
		//fmt.Println()
	}
}

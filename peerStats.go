package main

import (
	"net"
	"time"
)

func peerStatLoop() {
	for {
		// Sleep until its time to update
		time.Sleep(time.Duration(SystemConfig.DataSources.CjdnsPeers.Interval))

		// Get peer stats from cjdns
		results, err := Data.CjdnsConn.InterfaceController_peerStats()
		if err != nil {
			l.Noticeln(err)
			return
		}

		Data.Node.RateIn, Data.Node.RateOut = 0, 0
		Data.Node.BytesIn, Data.Node.BytesOut = 0, 0
		peerUpdate := make(map[string]Peer)
		// loop through the results and update peer statistics while checking
		// for any peers that are no longer there.
		for _, peer := range results {
			// Create the peer in the map if it's not there yet and calculate
			// the IPv6 address. If it is there, just get the IPv6 address.
			var IPv6 net.IP
			if _, ok := Data.Peers[peer.PublicKey.String()]; !ok {
				IPv6 = peer.PublicKey.IP()
				Data.Peers[peer.PublicKey.String()] = Peer{
					IPv6: IPv6,
				}
			} else {
				IPv6 = Data.Peers[peer.PublicKey.String()].IPv6
			}

			// Calculate upload and download rate
			newRateIn := float64(peer.BytesIn-Data.Peers[peer.PublicKey.String()].BytesIn) / (time.Since(Data.Peers[peer.PublicKey.String()].LastUpdate).Seconds())
			newRateOut := float64(peer.BytesOut-Data.Peers[peer.PublicKey.String()].BytesOut) / (time.Since(Data.Peers[peer.PublicKey.String()].LastUpdate).Seconds())

			// Convert the last packet received timestamp to a time.Time
			last := time.Unix(0, peer.Last*1000000)

			// Update the peer statistics
			peerUpdate[peer.PublicKey.String()] = Peer{
				PublicKey:   peer.PublicKey,
				IPv6:        Data.Peers[peer.PublicKey.String()].IPv6, // save the same IP
				LastUpdate:  time.Now(),
				Last:        last,
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

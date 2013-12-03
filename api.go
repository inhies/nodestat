package main

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
)

// Handle requests to /all/ by returning all stats
func allStatsHandler(w http.ResponseWriter, r *http.Request) {
	// First check to see if we should allow this request.
	auth, err := SystemConfig.Access.JSONApi.Authentication.method(r)
	if err != nil {
		http.Error(w, "500 Server Error", http.StatusInternalServerError)
		l.Errln(err)
		return
	}

	if !auth {
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		return
	}

	err, httpStatusCode := updateCjdnsStats()
	if err != nil {
		http.Error(w, "500 Server Error", httpStatusCode)
		l.Errln(err)
		return
	}

	// Render the json and send it
	err = sendJSON(w, r, Data)
	if err != nil {
		l.Errln(err)
		return
	}
}

// Handle requests to /node/ by returning just cjdns stats
func nodeStatsHandler(w http.ResponseWriter, r *http.Request) {
	// First check to see if we should allow this request.
	auth, err := SystemConfig.Access.JSONApi.Authentication.method(r)
	if err != nil {
		http.Error(w, "500 Server Error", http.StatusInternalServerError)
		l.Errln(err)
		return
	}

	if !auth {
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		return
	}

	err, httpStatusCode := updateCjdnsStats()
	if err != nil {
		http.Error(w, "500 Server Error", httpStatusCode)
		l.Errln(err)
		return
	}

	err = sendJSON(w, r, Data.Node)
	if err != nil {
		http.Error(w, "500 Server Error", http.StatusInternalServerError)
		l.Errln(err)
		return
	}
}

// Handle requests to /peers/ by returning just peer stats
func peerStatsHandler(w http.ResponseWriter, r *http.Request) {
	l.Debugln("Received request for peer data")
	// First check to see if we should allow this request.
	auth, err := SystemConfig.Access.JSONApi.Authentication.method(r)
	if err != nil {
		http.Error(w, "500 Server Error", http.StatusInternalServerError)
		l.Errln(err)
		return
	}

	if !auth {
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		return
	}

	// Render the json and send it
	err = sendJSON(w, r, Data.Peers)
	if err != nil {
		l.Errln(err)
		return
	}
}

func sendJSON(w http.ResponseWriter, r *http.Request, v interface{}) (err error) {

	// Render the json and send it
	jsonOut, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return
	}

        r.ParseForm()
        callback := []byte(r.Form.Get("callback"))
        cb := []byte(r.Form.Get("cb"))
        if SystemConfig.Access.JSONApi.EnableJSCallbacks {
		if validIP.MatchString(r.Header.Get("Referer")) {
			l.Infoln("Successful JS access attempt from ", r.Header.Get("Referer"))
			if(len(callback) > 0) {
				callback := append(callback, []byte("(")...)
				jsonOut = append(callback, jsonOut...)
				jsonOut = append(jsonOut, []byte(")")...)
			} else if(len(cb) > 0) {
				cb := append(cb, []byte("(")...)
				jsonOut = append(cb, jsonOut...)
				jsonOut = append(jsonOut, []byte(")")...)
			}
		} else {
			http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		}
        }

	w.Header().Set("Content-Length", strconv.Itoa(len(jsonOut)))
	w.Header().Set("Content-Type", "Text/JavaScript")
	w.Write(jsonOut)
	return
}

// Always allows access, effectively disabling authentication.
func nullAuth(r *http.Request) (authorized bool, err error) {
	host, port, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		l.Errln(err)
		return
	}
	l.Infoln("Successful access attempt from", host, port)
	return true, nil
}

// Only allows access from specific IP addresses.
func IPAuth(r *http.Request) (authorized bool, err error) {
	host, port, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		l.Errln(err)
		return
	}

	for _, ip := range SystemConfig.Access.JSONApi.Authentication.IP.Authorized {
		if host == ip {
			l.Infoln("Successful access attempt from", host, port)
			return true, nil
		}
	}

	l.Infoln("Failed access attempt from", host)
	return
}

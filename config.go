package main

import (
	"encoding/json"
	"fmt"
	"github.com/inhies/go-log"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	// Location of the default configuration file. By default it should be in
	// the same directory as the nodestat executable.
	DefaultConfigFile = "nodestat.conf"
)

// NodeStat configuration.
type configuration struct {

	// To add: .cjdnsadmin file location

	// Options for specific types of data
	DataSources struct {
		// Settings for getting peer info from cjdns
		CjdnsPeers struct {
			// How frequently should we query cjdns for peer updates
			Interval Duration
		}
	}

	// Settings for the various means you can use NodeStat.
	Access struct {

		// Settings for the HTTP JSON API.
		JSONApi struct {

			// Is this method enabled.
			Enabled bool

			// Protocol, address, and port to listen on.
			Addr string

			// Settings for user authentication.
			Authentication struct {

				// The auth method as read from the configuration file.
				RawMethod string `json:"Method"`

				// The function to use for authenticating requests. This is here to
				// allow easy expansion to different methods.
				method func(*http.Request) (bool, error)

				// For the IP address authentication method, a list of authorized
				// addresses.
				IP struct {
					Authorized []string
				}
			}
		}

		// Settings for the Munin plugin.
		Munin struct {
			IsTheMuninPluginDoneYet bool
		}
	}

	// Settings for the log.
	Log struct {
		// Format for the log timestamp.
		TimestampFormat string

		// Timezone to log in.
		TimestampTimezone string

		// Used internally for logging the correct timezone.
		timestampLocation *time.Location

		// Log level as read from the config file.
		RawLogLevel interface{} `json:"LogLevel"`

		// Parsed log level.
		LogLevel log.LogLevel

		// Include the log message level in the log.
		IncLogLevel bool
	}
	
	// Settings for the front-end.
	Web struct {
		// Set to true if they want to see a front end, false
		// if they only want the API.
		EnableFrontEnd bool
	}
}

// Creates a JSON configuration file. I find it's easier to do set the values
// in the struct and then create the JSON than it is to hand write the JSON
// and any changes to the structure of it.
func makeConfig() {
	c := SystemConfig

	// Update cjdns peer information every second
	c.DataSources.CjdnsPeers.Interval = Duration(1 * time.Second)

	// Enable the HTTP JSON API
	c.Access.JSONApi.Enabled = true

	// Listen on all interfaces
	c.Access.JSONApi.Addr = "[::]:8080"

	// Use IPv6 addresses as authentication
	c.Access.JSONApi.Authentication.method = IPAuth
	c.Access.JSONApi.Authentication.RawMethod = "IP"

	// No default approved IP's
	c.Access.JSONApi.Authentication.IP.Authorized = make([]string, 0)

	// Only allow connections from localhost
	c.Access.JSONApi.Authentication.IP.Authorized = append(
		c.Access.JSONApi.Authentication.IP.Authorized,
		"::1", "127.0.0.1")

	// Debug level loggint
	c.Log.LogLevel = log.DEBUG
	c.Log.RawLogLevel = c.Log.LogLevel.String()
	c.Log.TimestampFormat = time.StampMicro

	//IANA Standard http://en.wikipedia.org/wiki/List_of_tz_database_time_zones
	c.Log.TimestampTimezone = "UTC"

	// Set the timestamp location based on the specified log timezone
	var err error
	c.Log.timestampLocation, err = time.LoadLocation(c.Log.TimestampTimezone)
	if err != nil {
		println("Invalid timestamp timezone specified")
		return
	}

	jsonout, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		fmt.Println("Error making config:", err)
	}
	fmt.Println(string(jsonout))
}

// Reloads the configuration file. This is useful for changing the log level
// without restarting.
func reloadConfig() (err error) {
	err = parseConfig()
	if err != nil {
		return err
	}

	logLevel, err := log.ParseLevel(SystemConfig.Log.RawLogLevel)
	if err != nil {
		return err
	}

	l.Println("Config reloaded. Log level:", logLevel)
	SystemConfig.Log.LogLevel = logLevel
	l.Level = logLevel
	return
}

// Loads and parses the configuration file
func parseConfig() (err error) {
	rawFile, err := ioutil.ReadFile(DefaultConfigFile) //tUser.HomeDir + "/.file"
	if err != nil {
		return err
	}
	raw := rawFile

	err = json.Unmarshal(raw, &SystemConfig)
	if err != nil {
		return err
	}

	// Set the authentication function based on the text in the config file.
	switch strings.ToUpper(SystemConfig.Access.JSONApi.Authentication.RawMethod) {
	case "IP":
		SystemConfig.Access.JSONApi.Authentication.method = IPAuth
	case "NONE", "DISABLED", "", "OFF":
		SystemConfig.Access.JSONApi.Authentication.method = nullAuth
	default:
		return fmt.Errorf("Invalid authentication method specified")
	}
	return
}

// Thanks to the nodeatlas team for this trick.
type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	if b[0] != '"' {
		// If the duration is not a string, then consider it to be the
		// zero duration, so we do not have to set it.
		return nil
	}
	dur, err := time.ParseDuration(string(b[1 : len(b)-1]))
	if err != nil {
		return err
	}
	*d = *(*Duration)(&dur)
	return nil
}

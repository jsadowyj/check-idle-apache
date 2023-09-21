package main

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-plugin-sdk/sensu"
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	URL string
	Warning int
	Critical int
}

const IDLE_NULL = -1

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-plugins-idle-apache",
			Short:    "Monitors idle workers on apache status page",
			Keyspace: "sensu.io/plugins/sensu-plugins-idle-apache/config",
		},
	}

	options = []sensu.ConfigOption{
		&sensu.PluginConfigOption[string]{
			Path:      "url",
			Env:       "URL",
			Argument:  "url",
			Shorthand: "u",
			Default:   "http://127.0.0.1/server-status?auto",
			Usage:     "mod_status url (?auto query parameter is required!!)",
			Value:     &plugin.URL,
		},
		&sensu.PluginConfigOption[int]{
			Path:      "warning",
			Env:       "WARNING",
			Argument:  "warning",
			Shorthand: "w",
			Default:   IDLE_NULL,
			Usage:     "warning threshold",
			Value:     &plugin.Warning,
		},
		&sensu.PluginConfigOption[int]{
			Path:      "critical",
			Env:       "CRITICAL",
			Argument:  "critical",
			Shorthand: "c",
			Default:   0,
			Usage:     "critical threshold",
			Value:     &plugin.Critical,
		},
	}
)

func main() {
	useStdin := false
	fi, err := os.Stdin.Stat()
	if err != nil {
		fmt.Printf("Error check stdin: %v\n", err)
		panic(err)
	}
	//Check the Mode bitmask for Named Pipe to indicate stdin is connected
	if fi.Mode()&os.ModeNamedPipe != 0 {
		fmt.Println("using stdin")
		useStdin = true
	}

	check := sensu.NewCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, useStdin)
	check.Execute()
}

func checkArgs(event *corev2.Event) (int, error) {
	// return unknown status if URL is unable to be parsed
	u, err := url.ParseRequestURI(plugin.URL)
	// verifies the input URL is valid
	if err != nil {
		return sensu.CheckStateUnknown, fmt.Errorf("invalid url (input: %s)", plugin.URL)
	}
	// verifies the URL is machine readable
	// https://httpd.apache.org/docs/2.4/mod/mod_status.html
	if !u.Query().Has("auto") {
		return sensu.CheckStateUnknown, fmt.Errorf("missing ?auto query parameter (example: http://127.0.0.1/server-status?auto)")
	}
	return sensu.CheckStateOK, nil
}

func sendApacheRequest(u string) (*http.Response, error) {
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	return resp, nil
}	

func parseIdleWorkers(s string) (int, error) {
	idleStr := strings.Split(s, " ")[1]
	if i, err := strconv.Atoi(idleStr); err == nil {
		return i, nil 
	} 
	return IDLE_NULL, fmt.Errorf("unable to parse idle worker count")
}

func parseApacheResponse(r *http.Response) (int, error) {
	scanner := bufio.NewScanner(r.Body)
	for scanner.Scan() {
		line := strings.ToLower(scanner.Text())
    if strings.Contains(line, "idleworkers:"){
			idle, err := parseIdleWorkers(line)
			return idle, err
		}
	}
	return IDLE_NULL, fmt.Errorf("unable to parse http response")
}

func executeCheck(event *corev2.Event) (int, error) {
	res, err := sendApacheRequest(plugin.URL)
	defer res.Body.Close()
	

	if err != nil {
		return sensu.CheckStateUnknown, err
	}

	idle, err := parseApacheResponse(res)

	if err != nil {
		return sensu.CheckStateUnknown, err
	}
	
	switch {
		case idle <= plugin.Critical:
			fmt.Printf("CRITICAL - idle workers: %d\n", idle)
			return sensu.CheckStateCritical, nil
		case idle <= plugin.Warning:
			fmt.Printf("WARNING - idle workers: %d\n", idle)
			return sensu.CheckStateWarning, nil
		default:
			fmt.Printf("OK - idle workers: %d\n", idle)
			return sensu.CheckStateOK, nil
	}
}


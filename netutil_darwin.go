/*
Package netutil exposes the ability to interrogate and modify certain aspects of the system's network
configuration.

The darwin (OS X) implementation uses the networksetup and scutil command line programs.  Clients that attempt
to modify network configuration will be prompted for the user's credentials in order to run the networksetup
program.  You can work around this by running your program as root.  One pattern is to put your network configuration
logic into a utility program, which you then chown root:wheel and chmod 4755 to let anyone execute it with
root privileges.  The reason for isolating this in a utility program is that running as root is very dangerous, so you
want to limit the damage that the program can do.
*/

package netutil

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

type MacNetInterface struct {
	name string
}

const (
	networksetup = "/usr/sbin/networksetup"
	scutil       = "/usr/sbin/scutil"
)

func IsInternetReachable() (reachable bool, err error) {
	var output string
	if output, err = call(scutil, "-r", "s3.amazonaws.com"); err != nil {
		return
	}
	reachable = (strings.TrimSpace(output) == "Reachable")
	return
}

func ListInterfaces() (intfs NetInterfaces, err error) {
	var output string
	if output, err = call(networksetup, "-listallnetworkservices"); err != nil {
		return
	}
	names := strings.Split(strings.TrimSpace(output), "\n")
	intfs = make(NetInterfaces, len(names)-1)
	for i, name := range names {
		// Skip first line, which is an informational message
		if i > 0 {
			intfs[i-1] = &MacNetInterface{name}
		}
	}
	return
}

func (intf *MacNetInterface) Name() string {
	return intf.name
}

func (intf *MacNetInterface) EnableHTTPProxy(addr string) (err error) {
	splitAddr := strings.Split(addr, ":")
	if len(splitAddr) != 2 {
		err = fmt.Errorf("Please specify addr as host:port (e.g. 127.0.0.1:8080)")
		return
	}
	_, err = call(networksetup, "-setwebproxy", intf.name, splitAddr[0], splitAddr[1])
	return
}

func (intf *MacNetInterface) DisableHTTPProxy() (err error) {
	_, err = call(networksetup, "-setwebproxystate", intf.name, "off")
	return
}

func call(cmdName string, args ...string) (output string, err error) {
	// Set up our command
	cmd := exec.Command(cmdName, args...)

	// Run the command in a goroutine
	pOut := make(chan string)
	pErr := make(chan error)
	go func() {
		if bytes, err := cmd.Output(); err != nil {
			pErr <- err
		} else {
			pOut <- string(bytes)
		}
	}()

	// Wait for the command to finish, time out in 5 seconds
	select {
	case <-time.After(5 * time.Second):
		if err := cmd.Process.Kill(); err != nil {
			log.Printf("failed to time out process: %s", err)
		}
	case output = <-pOut:
		return
	case err = <-pErr:
		return
	}
	return
}

package netutil

import (
	"testing"
)

func TestIt(t *testing.T) {
	if reachable, err := IsInternetReachable(); err != nil {
		t.Fatalf("Unable to determine reachability: %s", err)
	} else if !reachable {
		t.Fatal("Internet not reachable")
	}
	if intfs, err := ListInterfaces(); err != nil {
		t.Fatalf("Unable to list interfaces: %s", err)
	} else {
		for _, intf := range intfs {
			if err := intf.EnableHTTPProxy("127.0.0.1:9050"); err != nil {
				t.Errorf("Unable to EnableHTTPProxy for intf '%s': %s", intf.Name(), err)
			} else {
				if err := intf.DisableHTTPProxy(); err != nil {
					t.Errorf("Unable to DisableHTTPProxy for intf '%s': %s", intf.Name(), err)
				}
			}
		}
	}
}

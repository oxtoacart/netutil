package netutil

/*
NetInterface encapsulates a system network interface.
*/
type NetInterface interface {
	/* Name() retrieves the name of this NetInterface. */
	Name() string

	/* EnableHTTPProxy turns on the use of the HTTP proxy at addr (host:port) for this NetInterface */
	EnableHTTPProxy(addr string) (err error)

	/* DisableHTTPProxy turns off the use of an HTTP proxy for this NetInterface */
	DisableHTTPProxy() (err error)
}

type NetInterfaces []NetInterface

func (intfs NetInterfaces) EnableHTTPProxy(addr string) (err error) {
	for _, intf := range intfs {
		if err = intf.EnableHTTPProxy(addr); err != nil {
			return
		}
	}
	return
}

func (intfs NetInterfaces) DisableHTTPProxy() (err error) {
	for _, intf := range intfs {
		if err = intf.DisableHTTPProxy(); err != nil {
			return
		}
	}
	return
}

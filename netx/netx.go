package netx

import (
	"net"
	"strconv"

	"github.com/projecteru2/pistage/errors"
)

var nicIPs = []string{}

func init() { // nolint
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}

	for _, ifaddr := range addrs {
		var ip net.IP
		switch typ := ifaddr.(type) {
		case *net.IPNet:
			ip = typ.IP
		case *net.IPAddr:
			ip = typ.IP
		default:
			continue
		}

		if ip.IsGlobalUnicast() {
			nicIPs = append(nicIPs, ip.String())
		}
	}
}

// GetLocalIP .
func GetLocalIP(network, laddr string) (string, error) {
	switch network {
	case "tcp", "tcp4", "tcp6":
		var tcpaddr, err = net.ResolveTCPAddr(network, laddr)
		if err != nil {
			return "", errors.Trace(err)
		}
		if tcpaddr.Port < 1 {
			return "", errors.Errorf("unexpectedly, resolve %s addr %s to %s", network, laddr, tcpaddr.String())
		}
		if len(nicIPs) < 1 {
			return "", errors.New("unknown IP addr")
		}
		return net.JoinHostPort(nicIPs[0], strconv.Itoa(tcpaddr.Port)), nil

	case "unix", "unixpacket":
		return laddr, nil

	default:
		return "", net.UnknownNetworkError(network)
	}
}

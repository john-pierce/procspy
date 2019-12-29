package procspy

// netstat reading.

import (
	"net"
	"strconv"
	"strings"
)

// parseDarwinNetstat parses netstat output. (Linux has ip:port, darwin
// ip.port. The 'Proto' column value also differs.)
func parseDarwinNetstat(out string) []Connection {
	//
	//  Active Internet connections
	//  Proto Recv-Q Send-Q  Local Address          Foreign Address        (state)
	//  tcp4       0      0  10.0.1.6.58287         1.2.3.4.443      		ESTABLISHED
	//
	res := []Connection{}
	for i, line := range strings.Split(out, "\n") {
		if i == 0 || i == 1 {
			// Skip header
			continue
		}

		// Fields are:
		fields := strings.Fields(line)
		if len(fields) != 6 {
			continue
		}

		t := Connection{
			Transport: "tcp",
		}

		locals := strings.Split(fields[3], ".")
		var localAddress = strings.Join(locals[:len(locals)-1], ".")
		var localPort = locals[len(locals)-1]
		if strings.Contains(localAddress, "*") {
			t.LocalAddress = net.IPv4(0, 0, 0, 0)
		} else {
			t.LocalAddress = net.ParseIP(localAddress)
		}
		if strings.Contains(localPort, "*") {
			localPort = "0"
		}
		p, err := strconv.Atoi(localPort)
		if err != nil {
			return nil
		}
		t.LocalPort = uint16(p)

		remotes := strings.Split(fields[4], ".")
		var remoteAddress = strings.Join(remotes[:len(remotes)-1], ".")
		var remotePort = remotes[len(remotes)-1]

		if strings.Contains(remoteAddress, "*") {
			t.RemoteAddress = net.IPv4(0, 0, 0, 0)
		} else {
			t.RemoteAddress = net.ParseIP(remoteAddress)
		}

		if strings.Contains(remotePort, "*") {
			remotePort = "0"
		}

		p, err = strconv.Atoi(remotePort)
		if err != nil {
			return nil
		}

		t.RemotePort = uint16(p)

		if t.State, err = tcpStateString(fields[5]); err != nil {
			t.State = 0
		}

		res = append(res, t)
	}

	return res
}

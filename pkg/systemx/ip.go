package systemx

import (
	"fmt"
	"net"
)

func GetLocalhostIp() (ip net.IP, err error) {
	var all []net.IP
	ip = net.IPv4(127, 0, 0, 1)

	var ift []net.Interface
	var addrs []net.Addr
	ift, err = net.Interfaces()
	if err != nil {
		return
	}
	for _, ifi := range ift {
		addrs, err = ifi.Addrs()
		if err != nil {
			continue
		}

		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ip = ipnet.IP
					if "eth0" == ifi.Name || "em1" == ifi.Name {
						all = append([]net.IP{ip}, all...)
					} else {
						all = append(all, ip)
					}

				}
			}
		}
	}
	if len(all) > 0 {
		ip = all[0]
	}

	return
}

func Ip2String(ip net.IP) string {
	return fmt.Sprintf("%d.%d.%d.%d", ip[12], ip[13], ip[14], ip[15])
}

func Ip2HexString(ip net.IP) string {
	return fmt.Sprintf("%02x%02x%02x%02x", ip[12], ip[13], ip[14], ip[15])
}

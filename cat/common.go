package cat

import (
	"strings"
)

type serverAddress struct {
	host string
	port int
}

func resolveServerAddresses(router string) (addresses []string) {
	for _, segment := range strings.Split(router, ";") {
		if len(segment) == 0 {
			continue
		}
		fragments := strings.Split(segment, ":")
		if len(fragments) != 2 {
			continue
		}

		addresses = append(addresses, segment)
	}

	return
}

func compareServerAddress(a, b *serverAddress) bool {
	if a == nil || b == nil {
		return false
	}
	if strings.Compare(a.host, b.host) == 0 {
		return a.port == b.port
	} else {
		return false
	}
}

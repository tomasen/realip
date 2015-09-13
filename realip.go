package realip

import (
	"log"
	"net"
	"net/http"
	"strings"
)

var lancidrs = []string{
	"127.0.0.1/8", "10.0.0.0/8", "169.254.0.0/16", "172.16.0.0/12", "192.168.0.0/16", "::1/128", "fc00::/7",
}

func isLocalAddress(addr string) bool {
	for _, it := range lancidrs {
		_, cidrnet, err := net.ParseCIDR(it)
		if err != nil {
			log.Println("ParseCIDR:", err) // assuming I did it right above
			return false
		}
		myaddr := net.ParseIP(addr)
		if cidrnet.Contains(myaddr) {
			return true
		}
	}
	return false
}

// Request.RemoteAddress contains port, which we want to remove i.e.:
// "[::1]:58292" => "[::1]"
func ipAddrFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}

// RealIP return client's real public IP address 
// from http request headers.
func RealIP(r *http.Request) string {
	hdr := r.Header
	hdrRealIP := hdr.Get("X-Real-Ip")
	hdrForwardedFor := hdr.Get("X-Forwarded-For")
	if len(hdrRealIP) == 0 && len(hdrForwardedFor) == 0 {
		return ipAddrFromRemoteAddr(r.RemoteAddr)
	}
	if hdrForwardedFor != "" {
		// X-Forwarded-For is potentially a list of addresses separated with ","
		parts := strings.Split(hdrForwardedFor, ",")
		for i, p := range parts {
			p = strings.TrimSpace(p)
			if isLocalAddress(p) {
				continue
			}
			parts[i] = p
		}
		if len(parts) > 0 {
			return parts[0]
		}
	}
	return hdrRealIP
}

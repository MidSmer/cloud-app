package route

import (
	"net"
	"net/http"
)

func GetIP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}

	ip, _, err := net.SplitHostPort(IPAddress)
	if err != nil {
		w.Write([]byte(""))
		return
	}

	userIP := net.ParseIP(ip)
	if userIP.To4() == nil {
		w.Write([]byte(""))
		return
	}

	w.Write([]byte(userIP.To4().String()))
}

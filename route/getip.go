package route

import (
	"errors"
	"net"
	"net/http"
)

func GetIP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("True-Client-Ip")
	}
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}

	getIP4 := func(addr string) (ip4 string, err error) {
		userIP := net.ParseIP(addr)
		if userIP.To4() == nil {
			return "", errors.New("not ip4")
		}

		return userIP.To4().String(), nil
	}

	ip, _, err := net.SplitHostPort(IPAddress)
	if err != nil {
		ip4, err := getIP4(IPAddress)
		if err != nil {
			w.Write([]byte(""))
			return
		}

		w.Write([]byte(ip4))
		return
	}

	ip4, err := getIP4(ip)
	if err != nil {
		w.Write([]byte(""))
		return
	}

	w.Write([]byte(ip4))
}

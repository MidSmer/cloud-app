package route

import (
	"net/http"
)

func GetHead(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	r.Header.Write(w)
}

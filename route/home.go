package route

import (
	"net/http"
)

func ServeHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/index.html")
}
package route

import (
	"net/http"
)

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/index.html")
}

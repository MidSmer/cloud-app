package route

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var (
	Logger *logrus.Logger
)

type Route struct {
	Name    string
	Methods []string
	Pattern string
	Prefix  bool
	Handler interface{}
	Func    bool
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	for _, route := range routes {
		var r *mux.Route

		if route.Prefix {
			r = router.PathPrefix(route.Pattern)
		} else {
			r = router.Path(route.Pattern)
		}

		if route.Func {
			r.HandlerFunc(route.Handler.(func(http.ResponseWriter, *http.Request)))
		} else {
			r.Handler(route.Handler.(http.Handler))
		}

		r.
			Methods(route.Methods...).
			Name(route.Name)
	}

	return router
}

var routes = Routes{
	Route{
		"Index",
		[]string{"GET"},
		"/",
		false,
		ServeIndex,
		true,
	},
	Route{
		"Static",
		[]string{"GET"},
		"/static/",
		true,
		http.StripPrefix(
			"/static/", http.FileServer(http.Dir("public/"))),
		false,
	},
	Route{
		"WS",
		[]string{"GET", "POST"},
		"/ws",
		false,
		WS,
		true,
	},
	Route{
		"GetIP",
		[]string{"GET"},
		"/get-ip",
		false,
		GetIP,
		true,
	},
	Route{
		"GetInfo",
		[]string{"GET"},
		"/get-info",
		false,
		GetInfo,
		true,
	},
	Route{
		"GetHead",
		[]string{"GET", "POST"},
		"/get-head",
		false,
		GetHead,
		true,
	},
}

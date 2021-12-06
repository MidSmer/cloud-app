package route

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
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
		ServeHome,
		true,
	},
	Route{
		"Index",
		[]string{"GET"},
		"/article/",
		true,
		ServeHome,
		true,
	},
	Route{
		"Index",
		[]string{"GET"},
		"/about",
		false,
		ServeHome,
		true,
	},
	Route{
		"Static",
		[]string{"GET"},
		"/static/",
		true,
		http.StripPrefix(
			"/static/", http.FileServer(http.Dir("public/static/"))),
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
		"API",
		[]string{"GET", "POST"},
		"/api/{name}",
		false,
		API,
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
		"CreateArticle",
		[]string{"POST"},
		"/create-article",
		false,
		CreateArticle,
		true,
	},
	Route{
		"UpdateArticle",
		[]string{"POST"},
		"/update-article",
		false,
		UpdateArticle,
		true,
	},
	Route{
		"FetchArticle",
		[]string{"POST"},
		"/fetch-article",
		false,
		FetchArticle,
		true,
	},
}

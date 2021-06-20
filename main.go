package main

import (
	"github.com/MidSmer/cloud-app/route"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var (
	logger *logrus.Logger
)

func main() {
	logger = logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	route.Logger = logger

	port := os.Getenv("PORT")

	if port == "" {
		logger.Error("did not get port!")
		port = "80"
	}

	router := mux.NewRouter()
	router.HandleFunc("/", route.ServeHome)
	router.HandleFunc("/ws", route.WS)
	router.HandleFunc("/get-ip", route.GetIP)
	router.HandleFunc("/get-info", route.GetInfo)
	http.Handle("/", router)
	httpErr := http.ListenAndServe(":"+port, nil)
	if httpErr != nil {
		logger.Fatal(httpErr)
	}
}

package main

import (
	"github.com/MidSmer/cloud-app/route"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var (
	log *logrus.Logger
)

func main() {
	log = logrus.New()

	log.SetLevel(logrus.DebugLevel)

	port := os.Getenv("PORT")

	if port == "" {
		log.Error("did not get port!")
		port = "80"
	}

	router := mux.NewRouter()
	router.HandleFunc("/", route.ServeHome)
	router.HandleFunc("/get-ip", route.GetIP)
	http.Handle("/", router)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

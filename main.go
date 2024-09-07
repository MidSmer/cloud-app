package main

import (
	"net/http"
	"os"

	"github.com/MidSmer/cloud-app/route"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	logger *logrus.Logger
)

func init() {
	viper.AutomaticEnv()
}

func main() {
	logger = logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	route.Logger = logger

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatal("read config failed: ", err)
	}

	port := os.Getenv("PORT")

	if port == "" {
		logger.Error("did not get port!")
		port = viper.GetString("Setting.Port")
	}

	router := route.NewRouter()
	httpErr := http.ListenAndServe(":"+port, router)
	if httpErr != nil {
		logger.Fatal(httpErr)
	}
}

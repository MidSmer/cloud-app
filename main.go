package main

import (
	"github.com/MidSmer/cloud-app/db"
	"github.com/MidSmer/cloud-app/route"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
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

	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		logger.Error("did not get database url!")
		dsn = viper.GetString("DB.DSN")
	}

	port := os.Getenv("PORT")

	if port == "" {
		logger.Error("did not get port!")
		port = viper.GetString("Setting.Port")
	}

	if err = db.Init(dsn); err != nil {
		logger.Fatal("db failed: ", err)
	}
	defer db.Close()

	router := route.NewRouter()
	httpErr := http.ListenAndServe(":"+port, router)
	if httpErr != nil {
		logger.Fatal(httpErr)
	}
}

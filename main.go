package main

import (
	"bytes"
	"github.com/MidSmer/cloud-app/route"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	v2rayCore "v2ray.com/core"
	_ "v2ray.com/core/main/distro/all"
	"v2ray.com/ext/tools/conf/serial"
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

	v2rayConfig := `
	{
		"log" : {
			"loglevel": "debug"
		},
		"inbounds": [
			{
				"port": 80,
				"protocol": "vmess",
				"settings": {
					"clients": [
						{
							"id": "129ef4a2-cdf4-4f00-bc78-c0211d373ad8",
							"alterId": 64
						}
					],
					"disableInsecureEncryption": true
				},
				"streamSettings": {
					"network": "ws",
					"wsSettings": {
					    "path": "/ws"
				    }
				}
			}
		],
		"outbounds": [
			{
				"protocol": "freedom"
			}
		]
	}
	`
	config, err := serial.LoadJSONConfig(bytes.NewReader([]byte(v2rayConfig)))
	if err != nil {
		logger.Error("failed to load v2ray config!", err)
		return
	}

	//server, v2rayErr := v2rayCore.New(config)
	_, v2rayErr := v2rayCore.New(config)
	if v2rayErr != nil {
		logger.Error("failed to create v2ray server!", v2rayErr)
		return
	}

	//server.Start()

	router := mux.NewRouter()
	router.HandleFunc("/", route.ServeHome)
	router.HandleFunc("/ws", route.WS)
	router.HandleFunc("/get-ip", route.GetIP)
	http.Handle("/", router)
	httpErr := http.ListenAndServe(":"+port, nil)
	if httpErr != nil {
		logger.Fatal(httpErr)
	}
}

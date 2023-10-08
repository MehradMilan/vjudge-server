package main

import (
	"encoding/json"
	"log/slog"
	"os"
	"vjudge/pkg/util"
)

var config struct {
	// Different judge types; equivalent to different types of containers
	Secret string `json:"secret"`
	// Where (ip:port) should we listen
	ListenAddress string `json:"listen"`
	// Which endpoint should we answer
	Endpoint string `json:"endpoint"`
}

// readConfig will read the config file
func readConfig(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		util.SlogFatal("cannot open config file", err)
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		util.SlogFatal("cannot decode config file", err)
	}
	slog.Info("read config file")
}

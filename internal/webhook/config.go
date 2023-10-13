package webhook

import (
	"encoding/json"
	"log/slog"
	"os"
	"vjudge/pkg/util"
)

var config struct {
	GitUsername   string `json:"username"`
	GitPassword   string `json:"password"`
	TmpDirectory  string `json:"tmpdir"`
	TestDirectory string `json:"testdir"`
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

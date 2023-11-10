package webhook

import (
	"encoding/json"
	"log/slog"
	"os"
	"vjudge/pkg/util"
)

type Homework struct {
	Questions []string `json:"questions"`
}

var config struct {
	GitUsername   string              `json:"username"`
	GitPassword   string              `json:"password"`
	TmpDirectory  string              `json:"tmpdir"`
	SRCDirectory  string              `json:"srcdir"`
	TestDirectory string              `json:"testdir"`
	JudgeName     string              `json:"judgeName"`
	JudgeEmail    string              `json:"judgeEmail"`
	Homeworks     map[string]Homework `json:"homeworks"`
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

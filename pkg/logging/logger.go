package logging

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"os"
)

const (
	LogsDir = "logs"
)

func NewLogger(name string) *zap.Logger {
	_ = os.Mkdir(LogsDir, os.ModePerm)

	rawJSON := []byte(fmt.Sprintf(`{
	  "level": "debug",
	  "encoding": "json",
	  "outputPaths": ["stdout", "./logs/%s"],
	  "errorOutputPaths": ["stderr"],
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
	  }
	}`, name))

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger := zap.Must(cfg.Build())
	defer logger.Sync()

	return logger
}

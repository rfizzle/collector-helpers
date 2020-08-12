package state

import (
	"errors"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

func InitCLIParams() {
	flag.String("state-path", "collector.state", "state file path")
}

func ValidateCLIParams() error {
	if viper.GetString("state-path") == "" {
		return errors.New("missing state file path param (--state-path)")
	}

	dir, _ := filepath.Split(viper.GetString("state-path"))

	if !pathExists(dir) {
		return errors.New("invalid state file path (--state-path)")
	}

	return nil
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}

	return false
}
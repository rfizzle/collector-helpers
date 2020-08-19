package config

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

func InitCLIParams() {
	flag.StringP("config", "c", "", "config file")
}

func CheckConfigParams() error {
	if viper.GetString("config") != "" {
		if !fileExists(viper.GetString("config")) {
			return fmt.Errorf("config file does not exist at: %v", viper.GetString("config"))
		}

		dir, file := filepath.Split(viper.GetString("config"))
		extWithDot := strings.ToLower(filepath.Ext(viper.GetString("config")))
		ext := strings.ReplaceAll(extWithDot, ".", "")

		supportedTypes := []string{"json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"}
		if !contains(supportedTypes, ext) {
			return fmt.Errorf("invalid config file type (%s) (supported: %s )", ext, strings.Join(supportedTypes[:], ", "))
		}

		fileName := strings.TrimSuffix(file, extWithDot)

		viper.SetConfigName(fileName)
		viper.SetConfigType(ext)
		viper.AddConfigPath(dir)

		err := viper.ReadInConfig() // Find and read the config file
		if err != nil { // Handle errors reading the config file
			return fmt.Errorf("Fatal error config file: %s \n", err)
		}
	}

	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
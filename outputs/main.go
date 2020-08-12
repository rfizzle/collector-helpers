package outputs

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

func InitCLIParams() {
	gcsInitParams()
	s3InitParams()
	stackdriverInitParams()
	httpInitParams()
	fileInitParams()
}

func ValidateCLIParams() error {
	if err := gcsValidateParams(); err != nil {
		return err
	}

	if err := s3ValidateParams(); err != nil {
		return err
	}

	if err := stackdriverValidateParams(); err != nil {
		return err
	}

	if err := httpValidateParams(); err != nil {
		return err
	}

	if err := fileValidateParams(); err != nil {
		return err
	}

	return nil
}

func WriteToOutputs(src, timestamp string) error {
	// Google Cloud Storage output
	if viper.GetBool("gcs") {
		gcsPath := fmt.Sprintf("%s_%s.log", viper.GetString("gcs-path"), timestamp)
		if err := gcsWrite(src, gcsPath, viper.GetString("gcs-bucket"), viper.GetString("gcs-credentials")); err != nil {
			return fmt.Errorf("unable to write to google cloud storage: %v", err)
		}
	}

	// Amazon S3 output
	if viper.GetBool("s3") {
		s3Path := fmt.Sprintf("%s_%s.log", viper.GetString("s3-path"), timestamp)
		if err := s3Write(src, s3Path, viper.GetString("s3-region"), viper.GetString("s3-bucket"), viper.GetString("s3-access-key-id"), viper.GetString("s3-secret-key"), viper.GetString("s3-storage-class")); err != nil {
			log.Fatalf("Unable to write to amazon s3: %v", err)
		}
	}

	// Stackdriver output
	if viper.GetBool("stackdriver") {
		if err := stackdriverWrite(src, viper.GetString("stackdriver-project"), viper.GetString("stackdriver-log-name"), viper.GetString("stackdriver-credentials"), "id.time"); err != nil {
			log.Fatalf("Unable to write to stackdriver: %v", err)
		}
	}

	// HTTP output
	if viper.GetBool("http") {
		if err := httpWrite(src, viper.GetString("http-url"), viper.GetString("http-auth"), viper.GetInt("http-max-items")); err != nil {
			log.Fatalf("Unable to write to HTTP: %v", err)
		}
	}

	// File output
	if viper.GetBool("file") {
		if size, err := fileWrite(src, viper.GetString("file-path"), viper.GetBool("file-rotate")); err != nil || size == 0 {
			return fmt.Errorf("unable to write %v bytes to file: %v", size, err)
		}
	}

	return nil
}

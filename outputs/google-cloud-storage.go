package outputs

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
	"io"
	"os"
)

// gcsInitParams initializes the required CLI params for google cloud storage output.
// Uses pflag to setup flag options.
func gcsInitParams() {
	flag.Bool("gcs", false, "enable google cloud storage output")
	flag.String("gcs-bucket", "", "google cloud storage bucket")
	flag.String("gcs-path", "", "google cloud storage file path")
	flag.String("gcs-credentials", "", "google cloud storage credential file")
}

// gcsValidateParams checks if the google cloud storage param has been set and validates related params.
// Uses viper to get parameters. Set in collectors as flags and environment variables.
func gcsValidateParams() error {
	if viper.GetBool("gcs") {
		if viper.GetString("gcs-bucket") == "" {
			return errors.New("missing google cloud storage bucket param (--gcs-bucket)")
		}
		if viper.GetString("gcs-path") == "" {
			return errors.New("missing google cloud storage output path param (--gcs-path)")
		}
		if !fileExists(viper.GetString("gcs-credentials")) {
			return errors.New("missing google cloud storage credential file (--gcs-credentials)")
		}
	}

	return nil
}

// gcsWrite takes the temporary storage file with results and copies it to google cloud storage.
func gcsWrite(src, dst, bucketName, credentialsFile string) error {
	// Setup context and storage client
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))

	// Handle client errors
	if err != nil {
		return err
	}

	// Open the source file
	source, err := os.Open(src)

	// Handle source file errors
	if err != nil {
		return err
	}

	// Define the google cloud storage file destination
	googleCloudStorageFile := client.Bucket(bucketName).Object(dst).NewWriter(ctx)

	// Upload the file
	if _, err = io.Copy(googleCloudStorageFile, source); err != nil {
		return err
	}

	// Handle google cloud storage file closure errors
	if err := googleCloudStorageFile.Close(); err != nil {
		return err
	}

	// Handle source file closure errors
	if err := source.Close(); err != nil {
		return err
	}

	// Handle storage client closure errors
	if err := client.Close(); err != nil {
		return err
	}

	// Output to debug
	log.Debugf("google cloud storage output written to : %s/%s", bucketName, dst)

	return nil
}

package outputs

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type TmpWriter struct {
	lock         sync.Mutex
	Fp           *os.File
	LastFilePath string
}

// Make a new TmpWriter. Return nil and error if error occurs during setup.
func NewTmpWriter() (*TmpWriter, error) {
	w := &TmpWriter{}

	// Open the file
	f, err := ioutil.TempFile("", randomStringWithLength(64))

	// Handle error
	if err != nil {
		return nil, err
	}

	//defer os.Remove(f.Name())

	// Set file pointer
	w.Fp = f

	return w, nil
}

// Perform the actual act of rotating and reopening file.
func (w *TmpWriter) Rotate() (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	// Close existing file if open
	if w.Fp != nil {
		w.LastFilePath = w.Fp.Name()
		err = w.Fp.Close()
		w.Fp = nil
		if err != nil {
			return err
		}
	}

	if viper.GetBool("verbose") {
		log.Printf("Temp file rotated \n")
	}

	// Create a file.
	w.Fp, err = ioutil.TempFile("", randomStringWithLength(64))
	return err
}

func (w *TmpWriter) WriteLog(message string) (err error) {
	if _, err := w.Fp.WriteString(message + "\n"); err != nil {
		return fmt.Errorf("Error writing string: %v\n", err)
	}

	return err
}

package outputs

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

type TmpWriter struct {
	lock         sync.Mutex
	Fp           *os.File
	LastFilePath string
}

// NewTmpWriter initializes a new TmpWriter object with an open temp file pointer reference.
func NewTmpWriter() (*TmpWriter, error) {
	w := &TmpWriter{}

	// Open the file
	err := w.Open()

	return w, err
}

// Open a new temp file and set pointer reference
func (w *TmpWriter) Open() (err error) {
	w.Fp, err = ioutil.TempFile("", randomStringWithLength(64))
	return err
}

// Close the currently open file and empty the pointer reference
func (w *TmpWriter) Close() (err error) {
	if w.Fp != nil {
		w.LastFilePath = w.Fp.Name()
		err = w.Fp.Close()
		w.Fp = nil
	}

	return err
}

// Perform the actual act of rotating and reopening file.
func (w *TmpWriter) Rotate() (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	// Close existing file if open
	if err = w.Close(); err != nil {
		return err
	}

	// Create a file.
	err = w.Open()
	return err
}

// WriteLog appends a message to the currently open temp file referenced in the pointer
func (w *TmpWriter) WriteLog(message string) (err error) {
	if _, err := w.Fp.WriteString(message + "\n"); err != nil {
		return fmt.Errorf("Error writing string: %v\n", err)
	}

	return err
}

// package fileutils contains some helper functions used in unit tests
package fileutils

import (
	"fmt"
	"io/ioutil"
	"os"
)

const fileNamePrefix = "tailer-test-"

func CreateTestFile() (*os.File, error) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), fileNamePrefix)
	if err != nil {
		return nil, err
	}
	return tmpFile, nil
}

func RemoveTestFile(f *os.File) error {
	if f == nil {
		return fmt.Errorf("cannot remove nil file")
	}
	return os.Remove(f.Name())
}

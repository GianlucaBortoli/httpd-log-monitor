// package testutils contains some helper functions used in unit tests
package testutils

import (
	"fmt"
	"io/ioutil"
	"os"
)

const fileName = "tailer-test-"

func CreateTestFile() (*os.File, error) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), fileName)
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

package testutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTestFile(t *testing.T) {
	f, err := CreateTestFile()
	assert.NoError(t, err)
	assert.NotNil(t, f)
}

func TestRemoveTestFile(t *testing.T) {
	f, createErr := CreateTestFile()
	assert.NoError(t, createErr)
	assert.NotNil(t, f)

	removeErr := RemoveTestFile(f)
	assert.NoError(t, removeErr)
}

func TestRemoveTestFile_NilFile(t *testing.T) {
	err := RemoveTestFile(nil)
	assert.Error(t, err)
}

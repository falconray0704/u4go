package u4go

import (
	"io/ioutil"
	"os"
	"testing"
)

// TempFile persists contents and returns the path and a clean func
func TempFile(t *testing.T, fileLocation, fileNamePrefix, contents string) (path string, clean func()) {
	content := []byte(contents)
	tmpfile, err := ioutil.TempFile(fileLocation, fileNamePrefix)
	if err != nil {
		t.Fatal("Unable to create tmpfile", err)
	}

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal("Unable to write tmpfile", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal("Unable to close tmpfile", err)
	}

	filePath := tmpfile.Name()
	return filePath, func() {
		_ = os.Remove(tmpfile.Name())
	}
}



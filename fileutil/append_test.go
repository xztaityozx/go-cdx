package fileutil

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppend(t *testing.T) {
	dir, err := ioutil.TempDir("", "go-cdx")
	assert.NoError(t, err)
	f := filepath.Join(dir, "append")
	err = Append(f, "element1")
	assert.NoError(t, err)

	out, err := ioutil.ReadFile(f)
	assert.NoError(t, err)

	assert.Equal(t, []byte("element1"), out)
	err = Append(f, "element2")
	assert.NoError(t, err)
	out, err = ioutil.ReadFile(f)
	assert.NoError(t, err)
	assert.Equal(t, []byte("element1element2"), out)

	os.RemoveAll(dir)
}

func TestAppendLine(t *testing.T) {

	newline := func() string {
		if runtime.GOOS == "windows" {
			return "\r\n"
		} else if runtime.GOOS == "darwin" {
			return "\r"
		} else {
			return "\n"
		}
	}()

	dir, err := ioutil.TempDir("", "go-cdx")
	assert.NoError(t, err)
	f := filepath.Join(dir, "append")
	err = AppendLine(f, "line1")
	assert.NoError(t, err)

	out, err := ioutil.ReadFile(f)
	assert.NoError(t, err)
	assert.Equal(t, []byte("line1"+newline), out)

	err = AppendLine(f, "line2")
	assert.NoError(t, err)

	out, err = ioutil.ReadFile(f)
	assert.NoError(t, err)
	assert.Equal(t, []byte("line1"+newline+"line2"+newline), out)

	os.RemoveAll(dir)
}

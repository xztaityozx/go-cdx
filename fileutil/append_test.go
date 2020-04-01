package fileutil

import (
	"io/ioutil"
	"os"
	"path/filepath"
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


package subcmd_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xztaityozx/go-cdx/subcmd"
)

func Test_Popd(t *testing.T) {
	o, err := subcmd.Popd()

	assert.Nil(t, err)
	assert.Equal(t, o, "popd")
}

func Test_Initialize(t *testing.T) {
	o, err := subcmd.Initialize()
	as := assert.New(t)

	as.Nil(err)

	t.Run("windows", func(t *testing.T) {
		if runtime.GOOS != "windows" {
			t.Skip()
		}

		as.Equal(`
function cdx() {
    begin{
        [System.Collections.Generic.List[string]]$paths=@{}
    }
    process{
        $paths.Add($_)
    }
    end{
        $command = "$($paths | go-cdx $args)";
        Invoke-Expression "$command"
    }
}`, o)
	})

	t.Run("unix", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip()
		}

		as.Equal(`function() cdx() {command="$(go-cdx $@)"; eval "${command}"}`, o)
	})
}

func Test_Add(t *testing.T) {
	baseDir := filepath.Join(os.TempDir(), "go-cdx-test")
	_ = os.MkdirAll(baseDir, 0755)

	as := assert.New(t)

	t.Run("ディレクトリには書き出せない", func(t *testing.T) {
		o, err := subcmd.Add(baseDir)
		as.Error(err)
		as.Equal("false", o)
	})

	bkmkFile := filepath.Join(baseDir, "bkmk")
	cwd, _ := os.Getwd()

	t.Run("ファイルがない状態でも書き出せる", func(t *testing.T) {
		o, err := subcmd.Add(bkmkFile)

		as.Nil(err)
		as.Equal("true", o)

		data, err := ioutil.ReadFile(bkmkFile)
		as.Equal([]byte(fmt.Sprintln(cwd)), data)
	})

	t.Run("ファイルがあったら追記する", func(t *testing.T) {
		o, err := subcmd.Add(bkmkFile)

		as.Nil(err)
		as.Equal("true", o)

		data, err := ioutil.ReadFile(bkmkFile)
		as.Equal([]byte(fmt.Sprintf("%s\n%s\n", cwd, cwd)), data)
	})

	_ = os.RemoveAll(baseDir)
}

func Test_GitRoot(t *testing.T) {
	as := assert.New(t)

	baseDir := filepath.Join(os.TempDir(), "go-cdx-test")
	_ = os.MkdirAll(baseDir, 0755)
	_ = os.Chdir(baseDir)

	err := exec.Command("git", "init").Run()
	as.Nil(err)

	d1 := filepath.Join(baseDir, "1")
	d2 := filepath.Join(d1, "2")

	_ = os.MkdirAll(d2, 0755)

	data := []struct {
		cwd    string
		expect string
	}{
		{cwd: baseDir, expect: "true"},
		{cwd: d1, expect: "pushd ../"},
		{cwd: d2, expect: "pushd ../../"},
	}

	for _, v := range data {
		_ = os.Chdir(v.cwd)
		o, err := subcmd.GitRoot()
		as.Nil(err)
		as.Equal(v.expect, o, fmt.Sprintf("cwd: %s", v.cwd))
	}

	_ = os.RemoveAll(baseDir)
}

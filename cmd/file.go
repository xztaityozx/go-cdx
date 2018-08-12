package cmd

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os/exec"
	"strings"
)

func writeInputFile(stream *[]Record) string {
	tmp, _ := ioutil.TempFile("", "cdxInput")
	defer tmp.Close()

	w := bufio.NewWriter(tmp)
	defer w.Flush()

	res := []string{}

	for _, v := range *stream {
		res = append(res, v.ToString())
	}

	if _, err := w.WriteString(strings.Join(res, "\n")); err != nil {
		Fatal(err)
	}

	return tmp.Name()
}

func getRecordsFromFile(f string) *[]Record {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		Fatal(err)
	}

	var res []Record
	if err := json.Unmarshal(b, &res); err != nil {
		Fatal(err)
	}

	return &res
}

func getRecordsFromCommand(command string) *[]Record {
	out, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		Fatal(err)
	}

	var rt []Record

	for i, v := range strings.Split(string(out), "\n") {
		rt = append(rt, Record{
			Number: i,
			Path:   v,
		})
	}

	return &rt

}

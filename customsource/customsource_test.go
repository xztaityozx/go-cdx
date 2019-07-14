package customsource

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomSource_String(t *testing.T) {
	data := []struct {
		expect string
		cs     CustomSource
	}{
		{expect: fmt.Sprintf("ten\tt\t10\techo 10"), cs: CustomSource{Name: "ten", Alias: 't', Command: "echo 10", BeginColumn: 10}},
	}

	for i, v := range data {
		t.Run(fmt.Sprintf("test-%d", i), func(t *testing.T) {
			actual := v.cs.String()
			assert.Equal(t, v.expect, actual)
		})
	}
}

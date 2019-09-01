package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCollection(t *testing.T) {
	cfg := Collection{
		CdxSource{Alias: "a",Command:"command"},
	}

	t.Run("ok", func(t *testing.T) {
		data := []struct{
			req string
			exp []string
		} {
			{req: "a", exp:[]string{"a"}},
			{req:"h", exp:[]string{"h"}},
			{req:"b", exp:[]string{"b"}},
			{req:"ah", exp:[]string{"a", "h"}},
			{req:"ahb", exp:[]string{"a","h", "b"}},
		}

		for _, v := range data {
			e, err := NewCollection("","",cfg, v.req)
			assert.NoError(t, err)
			var box []string
			for _, v := range e {
				box = append(box, v.Alias)
			}
			assert.Equal(t, box, v.exp)
		}
	})
	t.Run("ng", func(t *testing.T) {
		_, err := NewCollection("","",cfg,"x")
		assert.Error(t, err)
	})
}


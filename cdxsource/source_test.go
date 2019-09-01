package cdxsource

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCollection(t *testing.T) {
	cfg := Collection{
		Source{Alias:'a', Command:"command"},
	}

	t.Run("ok", func(t *testing.T) {
		data := []struct{
			req string
			exp []rune
		} {
			{req: "a", exp:[]rune{'a'}},
			{req:"h", exp:[]rune{'h'}},
			{req:"b", exp:[]rune{'b'}},
			{req:"ah", exp:[]rune{'a', 'h'}},
			{req:"ahb", exp:[]rune{'a','h', 'b'}},
		}

		for _, v := range data {
			e, err := NewCollection("","",cfg, v.req)
			assert.NoError(t, err)
			var box []rune
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


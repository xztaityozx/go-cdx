package customsource

import (
	"context"
	"fmt"
	"testing"

	"github.com/xztaityozx/go-cdx/environment"

	"github.com/b4b4r07/go-finder"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	source := []CustomSource{
		{Name: "A", SubName: "A"},
		{Name: "K", SubName: "K"},
		{Name: "B", SubName: "B"},
		{Name: "C", SubName: "C"},
		{Name: "L", SubName: "L"},
	}

	data := []struct {
		expect []string
		in     string
		isErr  bool
	}{
		{expect: []string{"A", "B", "C"}, in: "ABC", isErr: false},
		{expect: []string{"K"}, in: "K", isErr: false},
		{expect: []string{}, in: "X", isErr: true},
	}

	for _, v := range data {
		var actual []string
		res, err := New(v.in, source)
		for _, k := range res {
			actual = append(actual, k.Name)
		}

		if v.isErr {
			assert.Error(t, err)
		} else {
			assert.ElementsMatch(t, v.expect, actual)
		}

	}
}

func TestSourceCollection_Run(t *testing.T) {
	data := []struct {
		wants finder.Items
		sc    SourceCollection
		isErr bool
	}{
		{
			wants: finder.Items{
				{Key: "[seq-10]\t1", Value: []string{"1"}},
				{Key: "[seq-10]\t2", Value: []string{"2"}},
				{Key: "[seq-10]\t3", Value: []string{"3"}},
				{Key: "[seq-10]\t4", Value: []string{"4"}},
				{Key: "[seq-10]\t5", Value: []string{"5"}},
				{Key: "[seq-10]\t6", Value: []string{"6"}},
				{Key: "[seq-10]\t7", Value: []string{"7"}},
				{Key: "[seq-10]\t8", Value: []string{"8"}},
				{Key: "[seq-10]\t9", Value: []string{"9"}},
				{Key: "[seq-10]\t10", Value: []string{"10"}},
			},
			sc:    SourceCollection{{Name: "seq-10", BeginColumn: 0, Command: "seq 10"}},
			isErr: false,
		},
		{
			wants: finder.Items{
				{Key: "[seq-10|xargs-n2]\t1 2", Value: []string{"2"}},
				{Key: "[seq-10|xargs-n2]\t3 4", Value: []string{"4"}},
				{Key: "[seq-10|xargs-n2]\t5 6", Value: []string{"6"}},
				{Key: "[seq-10|xargs-n2]\t7 8", Value: []string{"8"}},
				{Key: "[seq-10|xargs-n2]\t9 10", Value: []string{"10"}},
			},
			sc:    SourceCollection{{Name: "seq-10|xargs-n2", BeginColumn: 1, Command: "seq 10|xargs -n2"}},
			isErr: false,
		},
		{
			wants: finder.Items{
				{Key: "[seq-10]\t1", Value: []string{"1"}},
				{Key: "[seq-10]\t2", Value: []string{"2"}},
				{Key: "[seq-10]\t3", Value: []string{"3"}},
				{Key: "[seq-10]\t4", Value: []string{"4"}},
				{Key: "[seq-10]\t5", Value: []string{"5"}},
				{Key: "[seq-10]\t6", Value: []string{"6"}},
				{Key: "[seq-10]\t7", Value: []string{"7"}},
				{Key: "[seq-10]\t8", Value: []string{"8"}},
				{Key: "[seq-10]\t9", Value: []string{"9"}},
				{Key: "[seq-10]\t10", Value: []string{"10"}},
				{Key: "[seq-10|xargs-n2]\t1 2", Value: []string{"2"}},
				{Key: "[seq-10|xargs-n2]\t3 4", Value: []string{"4"}},
				{Key: "[seq-10|xargs-n2]\t5 6", Value: []string{"6"}},
				{Key: "[seq-10|xargs-n2]\t7 8", Value: []string{"8"}},
				{Key: "[seq-10|xargs-n2]\t9 10", Value: []string{"10"}},
			},
			sc: SourceCollection{
				{Name: "seq-10", BeginColumn: 0, Command: "seq 10"},
				{Name: "seq-10|xargs-n2", BeginColumn: 1, Command: "seq 10|xargs -n2"},
			},
			isErr: false,
		},
		{
			sc:    SourceCollection{{Name: "err", Command: "exit 1", BeginColumn: 0}},
			isErr: true,
		},
	}

	for i, v := range data {
		t.Run(fmt.Sprintf("test-%d", i), func(t *testing.T) {
			as := assert.New(t)
			actual, err := v.sc.Run(context.Background(), environment.NewEnvironment())
			if v.isErr {
				as.Error(err)
			} else {
				as.NoError(err)
				as.ElementsMatch(v.wants, actual)
			}

		})
	}

}

package customsource

import (
	"context"
	"fmt"
	"testing"

	"github.com/b4b4r07/go-finder"
	"github.com/stretchr/testify/assert"
)

func TestSourceCollection_FindBeginColumn(t *testing.T) {

	data := []struct {
		wants int
		isErr bool
		sc    SourceCollection
		in    string
	}{
		{wants: 10, sc: SourceCollection{{Name: "ten", BeginColum: 10}, {Name: "nine", BeginColum: 9}}, in: "ten", isErr: false},
		{wants: 9, sc: SourceCollection{{Name: "ten", BeginColum: 10}, {Name: "nine", BeginColum: 9}}, in: "nine", isErr: false},
		{wants: -1, sc: SourceCollection{{Name: "ten", BeginColum: 10}, {Name: "nine", BeginColum: 9}}, in: "one", isErr: true},
	}

	for i, v := range data {
		t.Run(fmt.Sprintf("test-%d", i), func(t *testing.T) {
			res, err := v.sc.FindBeginColumn(v.in)
			assert.Equal(t, v.wants, res)
			if v.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
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
			sc:    SourceCollection{{Name: "seq-10", BeginColum: 0, Command: "seq 10"}},
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
			sc:    SourceCollection{{Name: "seq-10|xargs-n2", BeginColum: 1, Command: "seq 10|xargs -n2"}},
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
				{Name: "seq-10", BeginColum: 0, Command: "seq 10"},
				{Name: "seq-10|xargs-n2", BeginColum: 1, Command: "seq 10|xargs -n2"},
			},
			isErr: false,
		},
		{
			sc:    SourceCollection{{Name: "err", Command: "exit 1", BeginColum: 0}},
			isErr: true,
		},
	}

	for i, v := range data {
		t.Run(fmt.Sprintf("test-%d", i), func(t *testing.T) {
			as := assert.New(t)
			actual, err := v.sc.Run(context.Background())
			if v.isErr {
				as.Error(err)
			} else {
				as.NoError(err)
				as.ElementsMatch(v.wants, actual)
			}

		})
	}

}

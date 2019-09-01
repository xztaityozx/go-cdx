package cd

import (
	"context"
	"github.com/xztaityozx/go-cdx/config"
	"github.com/xztaityozx/go-cdx/cdxsource"
)

const command string = "pushd"

type Cd struct {
	cfg       config.Config
	candidate []string
}

func New(cfg config.Config, path []string) Cd {
	return Cd{
		cfg:       cfg,
		candidate: path,
	}
}

// Build build command for cdx
// param:
//   - ctx: context
// return:
//   - string: command string
func (cd Cd) Build(ctx context.Context) string {
	clc := cdxsource.NewCollection()
	return "echo NotImplement"
}


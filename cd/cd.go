package cd

import (
	"context"
	"github.com/xztaityozx/go-cdx/config"
)

const command string = "pushd"

type Cd struct {
	cfg config.Config
	dst []string
}

func New(cfg config.Config, path []string) Cd {
	return Cd{
		cfg: cfg,
		dst: path,
	}
}

// Build build command for cdx
// param:
//   - ctx: context
// return:
//   - string: command string
func (cd Cd) Build(ctx context.Context) string {
	return ""
}


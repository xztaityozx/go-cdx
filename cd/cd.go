package cd

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/xztaityozx/go-cdx/config"
	"os"
	"path/filepath"
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
//   - req: request
// return:
//   - string: command string
//   - error: error
func (cd Cd) Build(ctx context.Context, req string) (string, error) {

	if len(cd.candidate) == 0 && len(req) == 0 {
		return "true", nil
	}

	var path = ""
	if len(cd.candidate) >= 2 || len(req) != 0{
		clc, err := config.NewCollection(cd.cfg.HistoryFile, cd.cfg.BookmarkFile, cd.cfg.Source, req)
		if err != nil {
			return "", err
		}
		path, err = clc.Select(ctx, cd.cfg.FuzzyFinder, cd.candidate)
		if err != nil {
			return "", err
		}
	} else {
		path = cd.candidate[0]
	}

	path, _ = filepath.Abs(path)

	if _, err := os.Stat(path); err != nil {
		if cd.cfg.Make {
			logrus.Printf("%s not found. make ?", path)
			if y, e := cd.cfg.FuzzyFinder.YesNo(); e != nil || !y {
				return "", errors.Errorf("mkdir canceled")
			}
			if x := os.MkdirAll(path, 0755); x != nil {
				return "", errors.Errorf("failed mkdir")
			}
		} else {
			return "", errors.Errorf("make option not set")
		}
	}

	return func() string {
		if cd.cfg.NoOutput {
			return command + " \"" + path + "\" " + config.DevNull()
		} else {
			return command + "\"" + path + "\""
		}
	}(), nil
}


package tester

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/assay-it/assay-it/internal/config"
	"github.com/assay-it/assay-it/internal/gocc"
	"github.com/assay-it/assay-it/internal/printer"
)

type Evaler struct {
	cfg    *config.Config
	box    *gocc.Sandbox
	mod    *gocc.Package
	suites []*gocc.Suite
}

func NewEvaler(builddir, pkgName string, suites []string) (*Evaler, error) {
	cfg := &config.Config{
		Suites: suites,
	}

	if len(suites) == 0 {
		conf, err := config.NewFromPkg("")
		if err != nil {
			return nil, err
		}
		cfg.Suites = conf.Suites
	}

	if builddir == "" {
		builddir = filepath.Join(os.TempDir(), "assay-it")
	}

	// Config Sandbox
	box, err := gocc.NewSandbox(os.Stderr, builddir)
	if err != nil {
		return nil, fmt.Errorf("config sandbox failure %s: %w", builddir, err)
	}

	mod, err := gocc.NewPackage(box.Path, pkgName)
	if err != nil {
		return nil, fmt.Errorf("unable to config package %s: %w", pkgName, err)
	}

	return &Evaler{
		cfg:    cfg,
		box:    box,
		mod:    mod,
		suites: []*gocc.Suite{},
	}, nil
}

func (tt *Evaler) AnalyzeSourceCode(stderr *printer.Printer) error {
	tt.suites = []*gocc.Suite{}

	for _, suite := range tt.cfg.Suites {
		if stderr != nil {
			stderr.Info("use: %s\n", suite)
		}
		seq, err := tt.mod.CopyFrom(suite)
		if err != nil {
			return fmt.Errorf("unable to copy %s: %w", suite, err)
		}
		tt.suites = append(tt.suites, seq)
	}

	return nil
}

func (tt *Evaler) CreateRunner() error {
	if len(tt.suites) == 0 {
		return fmt.Errorf("no suites defined")
	}

	if err := tt.mod.CreateRunner(tt.suites); err != nil {
		return err
	}

	if err := tt.mod.CreateMod(); err != nil {
		return err
	}

	return nil
}

func (tt *Evaler) Compile() error {
	return tt.box.Compile(tt.mod)
}

func (tt *Evaler) Test(targetSUT string) ([]byte, error) {
	buf := bytes.Buffer{}
	err := tt.mod.Run(&buf, targetSUT)
	if err != nil {
		return buf.Bytes(), err
	}
	return buf.Bytes(), nil
}

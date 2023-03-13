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

type Tester struct {
	pkg    string
	cfg    *config.Config
	box    *gocc.Sandbox
	mod    *gocc.Module
	suites []*gocc.Suite
}

func NewTester(pkg string) (*Tester, error) {
	cfg, err := config.NewFromPkg(pkg)
	if err != nil {
		return nil, err
	}

	// Config Sandbox
	box, err := gocc.NewSandbox(os.Stderr, cfg.Runner)
	if err != nil {
		return nil, fmt.Errorf("config sandbox failure %s: %w", cfg.Runner, err)
	}

	// Config Package
	mod, err := gocc.NewModule("main", cfg.Module, box.Path)
	if err != nil {
		return nil, fmt.Errorf("unable to config module %s: %w", box.Path, err)
	}

	return &Tester{
		pkg:    pkg,
		cfg:    cfg,
		box:    box,
		mod:    mod,
		suites: []*gocc.Suite{},
	}, nil
}

func (tt *Tester) AnalyzeSourceCode(stderr *printer.Printer) error {
	tt.suites = []*gocc.Suite{}

	for _, suite := range tt.cfg.Suites {
		if stderr != nil {
			stderr.Info("use: %s\n", suite)
		}
		seq, err := tt.mod.CopyFrom(filepath.Join(tt.pkg, suite))
		if err != nil {
			return fmt.Errorf("unable to copy %s: %w", suite, err)
		}
		tt.suites = append(tt.suites, seq)
	}

	return nil
}

func (tt *Tester) CreateRunner() error {
	if len(tt.suites) == 0 {
		return fmt.Errorf("no suites defined for %s", tt.pkg)
	}

	return tt.mod.CreateRunner(tt.suites)
}

func (tt *Tester) Test(targetSUT string) ([]byte, error) {
	buf := bytes.Buffer{}
	err := tt.mod.Run(&buf, targetSUT)
	if err != nil {
		return buf.Bytes(), err
	}
	return buf.Bytes(), nil
}

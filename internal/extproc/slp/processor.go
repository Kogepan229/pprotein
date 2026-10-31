package slp

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/kaz/pprotein/internal/collect"
)

type (
	processor struct {
		confPath string
	}
)

// dbSubcommand returns the slp subcommand for the database selected by PPROTEIN_DB_TYPE.
func dbSubcommand() (string, error) {
	switch strings.ToLower(os.Getenv("PPROTEIN_DB_TYPE")) {
	case "", "mysql", "my", "mariadb":
		return "my", nil
	case "postgres", "postgresql", "pg":
		return "pg", nil
	default:
		return "", fmt.Errorf("unsupported PPROTEIN_DB_TYPE: %q", os.Getenv("PPROTEIN_DB_TYPE"))
	}
}

func (p *processor) Cacheable() bool {
	return true
}

func (p *processor) Process(snapshot *collect.Snapshot) (io.ReadCloser, error) {
	bodyPath, err := snapshot.BodyPath()
	if err != nil {
		return nil, fmt.Errorf("failed to find snapshot body: %w", err)
	}

	sub, err := dbSubcommand()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("slp", sub, "--config", p.confPath, "--output", "standard", "--format", "tsv", "--file", bodyPath)

	res, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("external process aborted: %w", err)
	}

	return io.NopCloser(bytes.NewBuffer(res)), nil
}

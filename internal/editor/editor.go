package editor

import (
	"context"
	"os"
	"os/exec"
)

func Detect(override string) string {
	if override != "" {
		return override
	}
	candidates := []string{os.Getenv("EDITOR"), "code", "nano", "vim", "vi"}
	for _, c := range candidates {
		if c == "" {
			continue
		}
		if p, err := exec.LookPath(c); err == nil {
			return p
		}
	}
	return "vi"
}
func Open(ctx context.Context, ed, path string) error {
	cmd := exec.CommandContext(ctx, ed, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

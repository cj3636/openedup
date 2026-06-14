package providers

import (
	"context"
	"github.com/Horryportier/openup/internal/entry"
	"os/exec"
	"strings"
)

func Zoxide(ctx context.Context, query string, limit int) []entry.Entry {
	if _, err := exec.LookPath("zoxide"); err != nil {
		return nil
	}
	args := []string{"query", "-l"}
	if query != "" {
		args = append(args, query)
	}
	out, err := exec.CommandContext(ctx, "zoxide", args...).Output()
	if err != nil {
		return nil
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	res := []entry.Entry{}
	for i, l := range lines {
		if l == "" {
			continue
		}
		if limit > 0 && i >= limit {
			break
		}
		res = append(res, entry.Entry{ID: l, Type: entry.FolderBrowse, Name: l, Description: "zoxide", Icon: "󰉋", Action: entry.ActionBrowseFolder, Metadata: map[string]string{"path": l}})
	}
	return res
}
func Systemd(ctx context.Context, limit int) []entry.Entry {
	if _, err := exec.LookPath("systemctl"); err != nil {
		return nil
	}
	out, err := exec.CommandContext(ctx, "systemctl", "list-units", "--type=service", "--state=running", "--no-legend", "--no-pager").Output()
	if err != nil {
		return nil
	}
	res := []entry.Entry{}
	for _, l := range strings.Split(string(out), "\n") {
		f := strings.Fields(l)
		if len(f) == 0 {
			continue
		}
		name := f[0]
		res = append(res, entry.Entry{ID: "svc:" + name, Type: entry.Command, Name: name, Description: "running systemd service", Icon: "●", Action: entry.ActionRunCommand, Metadata: map[string]string{"command": "systemctl status " + name}})
		if limit > 0 && len(res) >= limit {
			break
		}
	}
	return res
}

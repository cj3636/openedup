package history

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Horryportier/openup/internal/entry"
	"github.com/Horryportier/openup/internal/skate"
)

type Service struct{ Store *skate.Store }

func (s Service) Record(ctx context.Context, e entry.Entry) error {
	now := time.Now()
	stat := entry.Stat{Entry: e, Count: 1, FirstSeen: now, LastUsed: now}
	key := "history/" + e.ID
	var old entry.Stat
	if err := s.Store.GetJSON(ctx, key, &old); err == nil {
		stat = old
		stat.Count++
		stat.LastUsed = now
		stat.Entry = e
	}
	return s.Store.SetJSON(ctx, key, stat)
}
func (s Service) Entries(ctx context.Context) []entry.Entry {
	kv, _ := s.Store.Prefix(ctx, "history/")
	out := []entry.Entry{}
	for _, b := range kv {
		var st entry.Stat
		if json.Unmarshal(b, &st) == nil {
			st.Entry.Type = entry.History
			out = append(out, st.Entry)
		}
	}
	return out
}
func ShellCommands(limit int) []entry.Entry {
	home, _ := os.UserHomeDir()
	files := []string{filepath.Join(home, ".bash_history"), filepath.Join(home, ".zsh_history")}
	seen := map[string]bool{}
	out := []entry.Entry{}
	for _, p := range files {
		f, err := os.Open(p)
		if err != nil {
			continue
		}
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if line == "" || seen[line] {
				continue
			}
			seen[line] = true
			out = append(out, entry.Entry{ID: "cmd:" + line, Type: entry.Command, Name: line, Icon: "$", Action: entry.ActionRunCommand, Metadata: map[string]string{"command": line}})
			if limit > 0 && len(out) >= limit {
				f.Close()
				return out
			}
		}
		f.Close()
	}
	return out
}

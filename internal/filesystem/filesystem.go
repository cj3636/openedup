package filesystem

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Horryportier/openup/internal/entry"
)

func Entries(ctx context.Context, dir string, limit int) ([]entry.Entry, error) {
	if dir == "" {
		dir = "."
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	des, err := os.ReadDir(abs)
	if err != nil {
		return nil, err
	}
	sort.Slice(des, func(i, j int) bool {
		if des[i].IsDir() != des[j].IsDir() {
			return des[i].IsDir()
		}
		return strings.ToLower(des[i].Name()) < strings.ToLower(des[j].Name())
	})
	if limit <= 0 || limit > len(des) {
		limit = len(des)
	}
	out := make([]entry.Entry, 0, limit)
	for i := 0; i < limit; i++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		de := des[i]
		p := filepath.Join(abs, de.Name())
		typ, act, icon := entry.File, entry.ActionOpenFile, "󰈔"
		if de.IsDir() {
			typ, act, icon = entry.FolderView, entry.ActionViewFolder, "󰉋"
		}
		info := ""
		if fi, err := de.Info(); err == nil {
			info = fi.Mode().String()
		}
		out = append(out, entry.Entry{ID: p, Type: typ, Name: de.Name(), Description: p, Icon: icon, Action: act, Metadata: map[string]string{"path": p, "mode": info}})
	}
	return out, nil
}

func Preview(path string, max int) string {
	st, err := os.Stat(path)
	if err != nil {
		return err.Error()
	}
	if st.IsDir() {
		ents, err := os.ReadDir(path)
		if err != nil {
			return err.Error()
		}
		names := []string{}
		for i, e := range ents {
			if i >= 20 {
				break
			}
			suffix := ""
			if e.Type().IsDir() {
				suffix = "/"
			}
			names = append(names, e.Name()+suffix)
		}
		return strings.Join(names, "\n")
	}
	if st.Size() > 1<<20 {
		return "large file"
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return err.Error()
	}
	s := string(b)
	if max > 0 && len(s) > max {
		s = s[:max] + "\n…"
	}
	return s
}

func IsRegular(d fs.DirEntry) bool { t := d.Type(); return t.IsRegular() || t&fs.ModeSymlink != 0 }

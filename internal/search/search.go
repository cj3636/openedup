package search

import (
	"github.com/Horryportier/openup/internal/entry"
	"sort"
	"strings"
)

func Rank(q string, entries []entry.Entry, limit int) []entry.Entry {
	q = strings.ToLower(strings.TrimSpace(q))
	if q == "" {
		if limit > 0 && len(entries) > limit {
			return entries[:limit]
		}
		return entries
	}
	type scored struct {
		e entry.Entry
		s int
	}
	ss := []scored{}
	for _, e := range entries {
		hay := strings.ToLower(e.Name + " " + e.Description)
		if strings.Contains(hay, q) {
			score := 10
			if strings.HasPrefix(strings.ToLower(e.Name), q) {
				score += 20
			}
			if e.Type == entry.Favorite {
				score += 15
			}
			ss = append(ss, scored{e, score})
		}
	}
	sort.SliceStable(ss, func(i, j int) bool { return ss[i].s > ss[j].s })
	if limit <= 0 || limit > len(ss) {
		limit = len(ss)
	}
	out := make([]entry.Entry, limit)
	for i := 0; i < limit; i++ {
		out[i] = ss[i].e
		out[i].Type = entry.SearchResult
	}
	return out
}

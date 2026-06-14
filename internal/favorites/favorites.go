package favorites

import (
	"context"
	"encoding/json"
	"github.com/Horryportier/openup/internal/entry"
	"github.com/Horryportier/openup/internal/skate"
)

type Service struct{ Store *skate.Store }

func (s Service) Pin(ctx context.Context, e entry.Entry) error {
	st := entry.Stat{Entry: e, Count: 1, Pinned: true}
	return s.Store.SetJSON(ctx, "favorite/"+e.ID, st)
}
func (s Service) Entries(ctx context.Context) []entry.Entry {
	kv, _ := s.Store.Prefix(ctx, "favorite/")
	out := []entry.Entry{}
	for _, b := range kv {
		var st entry.Stat
		if json.Unmarshal(b, &st) == nil {
			st.Entry.Type = entry.Favorite
			out = append(out, st.Entry)
		}
	}
	return out
}

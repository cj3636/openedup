package skate

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Store struct {
	root string
	mu   sync.RWMutex
}

func Path(name string) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", "charm", "kv", name)
}

func Open(name string) (*Store, error) {
	p := Path(name)
	if err := os.MkdirAll(filepath.Join(p, "openedup-json"), 0o700); err != nil {
		return nil, err
	}
	return &Store{root: p}, nil
}
func (s *Store) Close() error { return nil }

func (s *Store) file(key string) string {
	sum := sha256.Sum256([]byte(key))
	return filepath.Join(s.root, "openedup-json", hex.EncodeToString(sum[:])+".json")
}
func (s *Store) index() string { return filepath.Join(s.root, "openedup-json", "index.json") }

func (s *Store) loadIndex() (map[string]string, error) {
	idx := map[string]string{}
	b, err := os.ReadFile(s.index())
	if errors.Is(err, os.ErrNotExist) {
		return idx, nil
	}
	if err != nil {
		return nil, err
	}
	return idx, json.Unmarshal(b, &idx)
}
func (s *Store) saveIndex(idx map[string]string) error {
	b, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.index(), b, 0o600)
}

func (s *Store) Set(ctx context.Context, key string, val []byte) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	idx, err := s.loadIndex()
	if err != nil {
		return err
	}
	idx[key] = filepath.Base(s.file(key))
	if err := os.WriteFile(s.file(key), val, 0o600); err != nil {
		return err
	}
	return s.saveIndex(idx)
}
func (s *Store) Get(ctx context.Context, key string) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return os.ReadFile(s.file(key))
}
func (s *Store) SetJSON(ctx context.Context, key string, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return s.Set(ctx, key, b)
}
func (s *Store) GetJSON(ctx context.Context, key string, v any) error {
	b, err := s.Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}
func (s *Store) Prefix(ctx context.Context, prefix string) (map[string][]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	idx, err := s.loadIndex()
	if err != nil {
		return nil, err
	}
	out := map[string][]byte{}
	for k := range idx {
		if strings.HasPrefix(k, prefix) {
			if b, err := os.ReadFile(s.file(k)); err == nil {
				out[k] = b
			}
		}
	}
	return out, nil
}

package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Command struct {
	Name, Description, Command string
	Confirm                    bool
}

type Config struct {
	Theme            string    `json:"theme"`
	Editor           string    `json:"editor"`
	FavoriteEditors  []string  `json:"favorite_editors"`
	FavoriteCommands []Command `json:"favorite_commands"`
	Preview          bool      `json:"preview"`
	FolderMode       string    `json:"folder_mode"`
	HistoryLimit     int       `json:"history_limit"`
	SearchLimit      int       `json:"search_limit"`
}

func Default() Config {
	return Config{Theme: "default", Preview: true, FolderMode: "view", HistoryLimit: 500, SearchLimit: 100}
}

func Path() string {
	if x := os.Getenv("XDG_CONFIG_HOME"); x != "" {
		return filepath.Join(x, "openedup", "config.json")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "openedup", "config.json")
}

func Load() (Config, error) {
	cfg := Default()
	b, err := os.ReadFile(Path())
	if errors.Is(err, os.ErrNotExist) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(b, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func Save(cfg Config) error {
	p := Path()
	if err := os.MkdirAll(filepath.Dir(p), 0700); err != nil {
		return err
	}
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, b, 0600)
}

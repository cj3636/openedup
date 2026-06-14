package app

import (
	"context"
	"github.com/Horryportier/openup/internal/config"
	"github.com/Horryportier/openup/internal/favorites"
	"github.com/Horryportier/openup/internal/history"
	"github.com/Horryportier/openup/internal/skate"
	"github.com/Horryportier/openup/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func Run(ctx context.Context, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	db, err := skate.Open("openedup")
	if err != nil {
		return err
	}
	defer db.Close()
	start := "."
	if len(args) > 0 {
		start = args[0]
	}
	m, err := ui.New(ctx, cfg, history.Service{Store: db}, favorites.Service{Store: db}, start)
	if err != nil {
		return err
	}
	return tea.NewProgram(m, tea.WithAltScreen()).Start()
}

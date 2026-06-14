package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Horryportier/openup/internal/app"
)

func main() {
	if err := app.Run(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

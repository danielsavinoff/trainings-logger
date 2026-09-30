package main

import (
	"context"
	"log"

	tea "charm.land/bubbletea/v2"
	"github.com/danielsavinoff/trainings-logger/internal/infra/persistence/sqlite"
	"github.com/danielsavinoff/trainings-logger/internal/infra/tui/starter"
)

func main() {
	ctx := context.Background()
	db, err := sqlite.Open(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := bootstrap(db)
	m, err := starter.NewModel(ctx, app)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		log.Fatal(err)
	}
}

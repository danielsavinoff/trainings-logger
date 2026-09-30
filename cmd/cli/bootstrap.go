package main

import (
	"database/sql"

	"github.com/danielsavinoff/trainings-logger/internal/app/command"
	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
	"github.com/danielsavinoff/trainings-logger/internal/app/query"
	"github.com/danielsavinoff/trainings-logger/internal/infra/persistence"
	"github.com/danielsavinoff/trainings-logger/internal/infra/persistence/sqlite"
	"github.com/danielsavinoff/trainings-logger/internal/infra/tui"
)

func bootstrap(db *sql.DB) tui.App {
	exercises := sqlite.NewExerciseRepository(db)
	sets := sqlite.NewSetRepository(db)
	preferences := sqlite.NewPreferencesRepository(db)
	transactions := persistence.NewTransactionRunner[ports.TransactionNoResult](
		db,
	)

	return tui.App{
		ListExercises: query.NewListExercisesHandler(
			exercises,
			persistence.NewTransactionRunner[query.ListExercisesResult](db),
		),
		ListSets: query.NewListSetsHandler(
			sets,
			persistence.NewTransactionRunner[query.ListSetsResult](db),
		),
		GetPreferences: query.NewGetPreferencesHandler(
			preferences,
			persistence.NewTransactionRunner[query.GetPreferencesResult](db),
		),
		CreateExercise: command.NewCreateExerciseHandler(
			exercises,
			persistence.NewTransactionRunner[command.CreateExerciseResult](db),
		),
		CreateSet: command.NewCreateSetHandler(
			preferences,
			sets,
			exercises,
			persistence.NewTransactionRunner[command.CreateSetResult](db),
		),
		UpdateExercise: command.NewUpdateExerciseHandler(
			exercises, transactions,
		),
		UpdateSet: command.NewUpdateSetHandler(sets, transactions),
		DeleteExercise: command.NewDeleteExerciseHandler(
			exercises, transactions,
		),
		DeleteSet: command.NewDeleteSetHandler(sets, transactions),
	}
}

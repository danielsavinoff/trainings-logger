package main

import (
	"database/sql"
	"os"

	"github.com/danielsavinoff/trainings-logger/internal/app/command"
	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
	"github.com/danielsavinoff/trainings-logger/internal/app/query"
	"github.com/danielsavinoff/trainings-logger/internal/infra/dropbox"
	"github.com/danielsavinoff/trainings-logger/internal/infra/persistence"
	"github.com/danielsavinoff/trainings-logger/internal/infra/persistence/sqlite"
	"github.com/danielsavinoff/trainings-logger/internal/infra/tui"
)

func bootstrap(db *sql.DB) tui.App {
	exercises := sqlite.NewExerciseRepository(db)
	sets := sqlite.NewSetRepository(db)
	preferences := sqlite.NewPreferencesRepository(db)
	identities := sqlite.NewIdentityRepository(db)
	credentials := sqlite.NewIdentityCredentialsRepository(db)
	transactions := persistence.NewTransactionRunner[ports.TransactionNoResult](
		db,
	)
	completeQuestionnaire := command.NewCompleteQuestionnaireHandler(
		preferences,
		identities,
		credentials,
		transactions,
	)
	linker := dropbox.NewLinker(os.Getenv("DROPBOX_APP_KEY"))
	stateExporter := sqlite.NewStateExporter(db)
	stateImporter := sqlite.NewStateImporter(db)
	backupStore := dropbox.NewStateStore(
		linker,
		identities,
		credentials,
	)
	saveRemoteState := command.NewSaveRemoteStateHandler(
		stateExporter,
		backupStore,
	)
	loadRemoteState := command.NewLoadRemoteStateHandler(
		stateImporter,
		backupStore,
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
		CompleteQuestionnaire: completeQuestionnaire.Handle,
		IdentityLinker:        linker,
		SaveRemoteState:       saveRemoteState.Handle,
		LoadRemoteState:       loadRemoteState.Handle,
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
			exercises,
			transactions,
		),
		UpdateSet: command.NewUpdateSetHandler(
			sets,
			transactions,
		),
		DeleteExercise: command.NewDeleteExerciseHandler(
			exercises,
			transactions,
		),
		DeleteSet: command.NewDeleteSetHandler(
			sets,
			transactions,
		),
	}
}

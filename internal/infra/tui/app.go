package tui

import (
	"context"

	"github.com/danielsavinoff/trainings-logger/internal/app/command"
	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
	"github.com/danielsavinoff/trainings-logger/internal/app/query"
)

type App struct {
	ListExercises         query.ListExercisesHandler
	ListSets              query.ListSetsHandler
	GetPreferences        query.GetPreferencesHandler
	CompleteQuestionnaire func(context.Context, command.CompleteQuestionnaire) error
	IdentityLinker        ports.IdentityLinker
	SaveRemoteState       func(context.Context) error
	LoadRemoteState       func(context.Context) error
	CreateExercise        command.CreateExerciseHandler
	CreateSet             command.CreateSetHandler
	UpdateExercise        command.UpdateExerciseHandler
	UpdateSet             command.UpdateSetHandler
	DeleteExercise        command.DeleteExerciseHandler
	DeleteSet             command.DeleteSetHandler
}

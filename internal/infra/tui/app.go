package tui

import (
	"github.com/danielsavinoff/trainings-logger/internal/app/command"
	"github.com/danielsavinoff/trainings-logger/internal/app/query"
)

// App supplies the application handlers used by the TUI.
type App struct {
	ListExercises  query.ListExercisesHandler
	ListSets       query.ListSetsHandler
	GetPreferences query.GetPreferencesHandler
	CreateExercise command.CreateExerciseHandler
	CreateSet      command.CreateSetHandler
	UpdateExercise command.UpdateExerciseHandler
	UpdateSet      command.UpdateSetHandler
	DeleteExercise command.DeleteExerciseHandler
	DeleteSet      command.DeleteSetHandler
}

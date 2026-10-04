// START_MODULE_CONTRACT
//   PURPOSE: Construct records state and retain help state.
//   SCOPE: Records model types, initial queries, and construction.
//   DEPENDS: M-TUI-HELP
//   LINKS: M-TUI-RECORDS, V-M-TUI-RECORDS
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   RecordsModel - retain records, dialogs, and help state
//   NewRecordsModel - load preferences and records
//   newRecordsModel - construct initial browsing state
//   Init - initialize the records model
// END_MODULE_MAP

package views

import (
	"context"

	"charm.land/bubbles/v2/paginator"
	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"

	"github.com/danielsavinoff/trainings-logger/internal/app/query"
	"github.com/danielsavinoff/trainings-logger/internal/domain/shared"
	"github.com/danielsavinoff/trainings-logger/internal/infra/tui"
	"github.com/danielsavinoff/trainings-logger/internal/infra/tui/components"
)

// START_BLOCK_RECORDS_MODEL
type panel int

const (
	exercisesPanel panel = iota
	setsPanel
)

type RecordsModel struct {
	keys           recordsKeyMap
	store          *recordStore
	weightUnit     shared.WeightUnit
	err            string
	exercises      []exerciseRecord
	setsExerciseID string
	table          table.Model
	timeline       viewport.Model
	exercisePager  paginator.Model
	setPager       paginator.Model
	setPageStarts  []int
	activePanel    panel
	goPending      bool
	setCursor      int
	setLines       []int
	width          int
	height         int
	dialog         *recordDialog
	help           helpModel
	showingHelp    bool
}

// NewRecordsModel loads the initial records and creates the terminal UI.
func NewRecordsModel(ctx context.Context, app tui.App) (RecordsModel, error) {
	store := &recordStore{App: app}
	prefs, err := store.GetPreferences.Handle(
		ctx,
		query.GetPreferences{},
	)
	if err != nil {
		return RecordsModel{}, err
	}

	data, err := store.ListExercises.Handle(
		ctx,
		query.ListExercises{},
	)
	if err != nil {
		return RecordsModel{}, err
	}

	exercises := make(
		[]exerciseRecord,
		0,
		len(data.Result),
	)
	for _, e := range data.Result {
		exercises = append(
			exercises,
			exerciseRecord{
				id:              e.ID,
				name:            e.Name,
				bodyWeightShare: e.BodyWeightShare,
			},
		)
	}

	m := newRecordsModel(exercises)
	m.store, m.weightUnit = store, prefs.Result.WeightUnit
	if err := m.loadSelectedSets(ctx); err != nil {
		return RecordsModel{}, err
	}

	m.refreshTimeline()

	return m, nil
}

func newRecordsModel(exercises []exerciseRecord) RecordsModel {
	m := RecordsModel{
		keys:          defaultRecordsKeyMap(),
		weightUnit:    shared.WeightUnitKg,
		exercises:     exercises,
		table:         newExerciseTable(exercises),
		timeline:      viewport.New(),
		exercisePager: components.NewPaginator(),
		setPager:      components.NewPaginator(),
		width:         100,
		height:        24,
	}

	return m.resize()
}

func (m RecordsModel) Init() tea.Cmd {
	return nil
}

var _ tea.Model = RecordsModel{}

// END_BLOCK_RECORDS_MODEL

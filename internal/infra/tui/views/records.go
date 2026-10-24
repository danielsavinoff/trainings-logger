// START_MODULE_CONTRACT
//   PURPOSE: Construct workout, catalog, membership, and set state.
//   SCOPE: Independent roots, domain-backed rows, startup setup, and transfer ownership.
//   DEPENDS: M-APP-SESSIONS, M-APP-PERFORMANCES, M-DOMAIN-SESSION, M-DOMAIN-PERFORMANCE, M-DOMAIN-EXERCISE, M-DOMAIN-SET
//   LINKS: M-TUI-RECORDS, V-M-TUI-RECORDS
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   RecordsModel - retain domain data, records focus, and catalog visibility
//   NewRecordsModel - load preferences and workout records
//   newRecordsModel - construct records/catalog tables and default session focus
//   Init - focus startup setup when required
// END_MODULE_MAP

package views

import (
	"context"

	"charm.land/bubbles/v2/paginator"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"

	"github.com/danielsavinoff/trainings-logger/internal/app/query"
	"github.com/danielsavinoff/trainings-logger/internal/domain/session"
	"github.com/danielsavinoff/trainings-logger/internal/domain/session/exercise"
	"github.com/danielsavinoff/trainings-logger/internal/domain/session/exercise/performance"
	"github.com/danielsavinoff/trainings-logger/internal/domain/session/exercise/set"
	"github.com/danielsavinoff/trainings-logger/internal/domain/shared"
	"github.com/danielsavinoff/trainings-logger/internal/infra/tui"
	"github.com/danielsavinoff/trainings-logger/internal/infra/tui/components"
)

// START_BLOCK_RECORDS_MODEL
type panel int

const (
	sessionsPanel panel = iota
	performancesPanel
	setsPanel
	catalogPanel
)

type RecordsModel struct {
	keys                  recordsKeyMap
	store                 *tui.App
	weightUnit            shared.WeightUnit
	err                   string
	catalog               []exercise.Exercise
	sessions              []session.Session
	performances          []performance.Performance
	sets                  []set.Set
	performancesSessionID string
	setsPerformanceID     string
	sessionTable          table.Model
	catalogTable          table.Model
	performanceTable      table.Model
	setTable              table.Model
	sessionPager          paginator.Model
	catalogPager          paginator.Model
	performancePager      paginator.Model
	setPager              paginator.Model
	activePanel           panel
	goPending             bool
	width                 int
	height                int
	dialog                *recordDialog
	help                  helpModel
	showingHelp           bool
	showingCatalog        bool
	catalogRecordsError   string
	questionnaire         *questionnaireModel
	transfer              *recordTransfer
	transferStatus        string
}

func NewRecordsModel(ctx context.Context, app tui.App) (RecordsModel, error) {
	prefs, err := app.GetPreferences.Handle(
		ctx,
		query.GetPreferences{},
	)
	if err != nil {
		return RecordsModel{}, err
	}

	m := newRecordsModel(nil)
	m.store = &app
	m.weightUnit = prefs.Result.WeightUnit
	if err := m.reloadRecords(
		ctx,
		"",
		"",
		"",
		"",
	); err != nil {
		return RecordsModel{}, err
	}

	if !prefs.Exists {
		m.openQuestionnaire()
	}

	return m, nil
}

func newRecordsModel(catalog []exercise.Exercise) RecordsModel {
	m := RecordsModel{
		keys:       defaultRecordsKeyMap(),
		weightUnit: shared.WeightUnitKg,
		catalog:    catalog,
		sessionTable: newRecordsTable(
			[]table.Column{{Title: "Index"}},
			nil,
		),
		catalogTable: newRecordsTable(
			[]table.Column{{Title: "Name"}, {Title: "Body Weight %"}},
			catalogRows(catalog),
		),
		performanceTable: newRecordsTable(
			[]table.Column{{Title: "Name"}, {Title: "Body Weight %"}},
			nil,
		),
		setTable: newRecordsTable(
			[]table.Column{{Title: "Reps"}, {Title: "Exercise Weight"}, {Title: "Rest (seconds)"}},
			nil,
		),
		sessionPager:     components.NewPaginator(),
		catalogPager:     components.NewPaginator(),
		performancePager: components.NewPaginator(),
		setPager:         components.NewPaginator(),
		width:            120,
		height:           30,
	}
	m.focusPanel(sessionsPanel)
	return m.resize()
}

func (m RecordsModel) Init() tea.Cmd {
	if m.questionnaire != nil {
		return m.questionnaire.input.Focus()
	}

	return nil
}

var _ tea.Model = RecordsModel{}

// END_BLOCK_RECORDS_MODEL

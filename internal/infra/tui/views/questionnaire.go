// START_MODULE_CONTRACT
//   PURPOSE: Collect optional first-start and reset answers in a full-screen view.
//   SCOPE: Staged answers, numeric/code input, asynchronous account linking, and return to records.
//   DEPENDS: M-APP-QUESTIONNAIRE, M-DOMAIN-SHARED-WEIGHT
//   LINKS: M-TUI-QUESTIONNAIRE, V-M-TUI-QUESTIONNAIRE, M-TUI-RECORDS
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   questionnaireModel - retain pending answers, input, and authorization context
//   questionnaireStartedMsg - deliver an authorization session to its owning view
//   questionnaireSavedMsg - deliver an exchange/persistence result to its owning view
//   openQuestionnaire - start an empty questionnaire without changing stored settings
//   closeQuestionnaire - cancel pending work and restore records
//   updateQuestionnaire - validate answers and route asynchronous completion
//   startQuestionnaireAuthorization - start or retry a fresh browser/code session
//   questionnaireURLContent - expose the complete scrollable authorization URL
//   saveQuestionnaire - save explicit answers, or return without writes when skipped
//   questionnaireView - render a bounded full-screen frame and native input cursor
// END_MODULE_MAP

package views

import (
	"context"
	"strings"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/danielsavinoff/trainings-logger/internal/app/command"
	"github.com/danielsavinoff/trainings-logger/internal/app/ports"
	"github.com/danielsavinoff/trainings-logger/internal/domain/shared"
	"github.com/danielsavinoff/trainings-logger/internal/domain/user/identity"
	"github.com/danielsavinoff/trainings-logger/internal/infra/tui/components"
)

// START_BLOCK_QUESTIONNAIRE_STATE
type questionnaireModel struct {
	step          int
	input         textinput.Model
	unit          *shared.WeightUnit
	authorization *ports.IdentityAuthorization
	linked        *ports.LinkedIdentity
	ctx           context.Context
	cancel        context.CancelFunc
	busy          bool
	status        string
	showURL       bool
	urlOffset     int
}

type questionnaireStartedMsg struct {
	owner         *questionnaireModel
	authorization ports.IdentityAuthorization
	err           error
}

type questionnaireSavedMsg struct {
	owner  *questionnaireModel
	linked *ports.LinkedIdentity
	err    error
}

func (m *RecordsModel) openQuestionnaire() tea.Cmd {
	if m.questionnaire != nil {
		m.questionnaire.cancel()
	}
	ctx, cancel := context.WithCancel(context.Background())
	input := textinput.New()
	input.SetVirtualCursor(false)
	input.Prompt = "› "
	input.Placeholder = "1 or 2"
	input.CharLimit = 1
	styles := input.Styles()
	styles.Focused.Prompt = lipgloss.NewStyle().Foreground(components.ActivePanelColor)
	styles.Cursor.Color = components.ActivePanelColor
	input.SetStyles(styles)
	m.questionnaire = &questionnaireModel{
		input:  input,
		ctx:    ctx,
		cancel: cancel,
	}

	m.goPending = false

	return m.questionnaire.input.Focus()
}

func (m *RecordsModel) closeQuestionnaire() {
	m.questionnaire.cancel()
	m.questionnaire = nil
}

func (m *RecordsModel) updateQuestionnaire(msg tea.Msg) (tea.Model, tea.Cmd) {
	q := m.questionnaire
	switch msg := msg.(type) {
	case questionnaireStartedMsg:
		if msg.owner != q {
			return *m, nil
		}

		q.busy = false
		if msg.err != nil {
			q.status = msg.err.Error()
			q.input.SetValue("")
			return *m, nil
		}

		q.authorization = &msg.authorization
		q.step = 2
		q.status = ""
		if !msg.authorization.BrowserOpened {
			q.status = "Browser could not open. Press ctrl+o for the sign-in URL."
		}
		q.input.SetValue("")
		q.input.CharLimit = 4096
		q.input.Placeholder = "Authorization code"
		q.input.EchoMode = textinput.EchoPassword
		return *m, nil

	case questionnaireSavedMsg:
		if msg.owner != q {
			return *m, nil
		}

		q.busy = false
		q.linked = msg.linked
		if msg.err != nil {
			q.status = msg.err.Error()
			if msg.linked == nil {
				q.input.SetValue("")
			}

			return *m, nil
		}

		if q.unit != nil {
			m.weightUnit = *q.unit
			if err := m.refreshSetRows(); err != nil {
				q.status = err.Error()
				return *m, nil
			}
			m.syncPagers()
		}
		m.closeQuestionnaire()
		return *m, nil

	case tea.KeyPressMsg:
		if q.showURL {
			switch msg.String() {
			case "esc", "ctrl+o":
				q.showURL = false
			case "j", "down", "k", "up", "pgdown", "pgup":
				width := max(
					1,
					m.width-components.PanelSpacing.GetHorizontalFrameSize()-components.BaseStyle.GetHorizontalFrameSize(),
				)
				height := max(
					1,
					m.height-components.PanelSpacing.GetVerticalFrameSize()-components.BaseStyle.GetVerticalFrameSize(),
				)
				limit := max(
					0,
					lipgloss.Height(m.questionnaireURLContent(width))-height,
				)
				delta := 1
				if msg.String() == "pgdown" || msg.String() == "pgup" {
					delta = height
				}

				if msg.String() == "k" || msg.String() == "up" || msg.String() == "pgup" {
					delta = -delta
				}

				q.urlOffset = min(
					limit,
					max(
						0,
						min(
							limit,
							q.urlOffset,
						)+delta,
					),
				)
			}

			return *m, nil
		}

		if msg.String() == "esc" {
			m.closeQuestionnaire()
			return *m, nil
		}

		if q.busy {
			return *m, nil
		}

		if q.step == 2 && msg.String() == "ctrl+r" {
			cmd := m.startQuestionnaireAuthorization(q)
			return *m, cmd
		}

		if q.step == 2 && msg.String() == "ctrl+o" && q.authorization != nil {
			q.showURL = true
			q.urlOffset = 0
			return *m, nil
		}

		if msg.String() != "enter" {
			break
		}

		answer := strings.TrimSpace(q.input.Value())
		switch q.step {
		case 0:
			if answer != "" && answer != "1" && answer != "2" {
				q.status = "Choose 1 or 2, or clear the input to skip."
				return *m, nil
			}

			if answer != "" {
				unit := shared.WeightUnitKg
				if answer == "2" {
					unit = shared.WeightUnitLb
				}

				q.unit = &unit
			}

			q.step = 1
			q.status = ""
			q.input.SetValue("")
			q.input.Placeholder = "1 or Enter to skip"
			return *m, nil

		case 1:
			if answer == "" {
				return *m, m.saveQuestionnaire(q)
			}

			if answer != "1" {
				q.status = "Choose 1, or clear the input to skip."
				return *m, nil
			}

			cmd := m.startQuestionnaireAuthorization(q)
			return *m, cmd

		case 2:
			return *m, m.saveQuestionnaire(q)
		}
	}

	if q.busy {
		return *m, nil
	}

	var cmd tea.Cmd
	q.input, cmd = q.input.Update(msg)

	return *m, cmd
}

func (m *RecordsModel) startQuestionnaireAuthorization(q *questionnaireModel) tea.Cmd {
	if m.store == nil || m.store.IdentityLinker == nil {
		q.status = "Dropbox sign-in is not configured. Clear the input to skip."
		return nil
	}
	q.cancel()
	next := *q
	next.ctx, next.cancel = context.WithCancel(context.Background())
	next.authorization, next.linked = nil, nil
	next.step = 1
	next.showURL, next.urlOffset = false, 0
	next.input.SetValue("")
	next.input.CharLimit = 1
	next.input.EchoMode = textinput.EchoNormal
	next.input.Placeholder = "1 or Enter to skip"
	next.busy = true
	next.status = "Opening Dropbox sign-in…"
	m.questionnaire = &next
	linker := m.store.IdentityLinker
	return func() tea.Msg {
		authorization, err := linker.Begin(
			next.ctx,
			identity.ProviderDropbox,
		)
		return questionnaireStartedMsg{
			owner:         &next,
			authorization: authorization,
			err:           err,
		}
	}
}

func (m *RecordsModel) saveQuestionnaire(q *questionnaireModel) tea.Cmd {
	code := strings.TrimSpace(q.input.Value())
	if q.step != 2 || code == "" {
		q.linked = nil
	}

	if q.unit == nil && (q.step != 2 || code == "") {
		m.closeQuestionnaire()
		return nil
	}

	if m.store == nil || m.store.CompleteQuestionnaire == nil {
		q.status = "Setup could not be saved. Press Esc to cancel."
		return nil
	}

	q.busy = true
	q.status = "Saving setup…"
	save := m.store.CompleteQuestionnaire
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(
			q.ctx,
			2*time.Minute,
		)
		defer cancel()
		linked := q.linked
		if q.step == 2 && code != "" && linked == nil {
			value, err := q.authorization.Complete(
				ctx,
				code,
			)
			if err != nil {
				return questionnaireSavedMsg{
					owner: q,
					err:   err,
				}
			}

			linked = &value
		}

		err := save(
			ctx,
			command.CompleteQuestionnaire{
				WeightUnit:     q.unit,
				LinkedIdentity: linked,
			},
		)

		return questionnaireSavedMsg{
			owner:  q,
			linked: linked,
			err:    err,
		}
	}
}

// END_BLOCK_QUESTIONNAIRE_STATE

// START_BLOCK_QUESTIONNAIRE_VIEW
func (m RecordsModel) questionnaireView() tea.View {
	v := tea.NewView("")
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	style := components.BaseStyle.BorderForeground(components.ActivePanelColor)
	panelWidth := m.width - components.PanelSpacing.GetHorizontalFrameSize()
	panelHeight := m.height - components.PanelSpacing.GetVerticalFrameSize()
	contentWidth := panelWidth - style.GetHorizontalFrameSize()
	contentHeight := panelHeight - style.GetVerticalFrameSize()
	if contentWidth <= 0 || contentHeight <= 0 {
		return v
	}

	q := m.questionnaire
	if q.showURL && q.authorization != nil {
		lines := strings.Split(
			m.questionnaireURLContent(contentWidth),
			"\n",
		)
		start := min(
			max(
				0,
				q.urlOffset,
			),
			max(
				0,
				len(lines)-contentHeight,
			),
		)
		body := strings.Join(
			lines[start:min(
				len(lines),
				start+contentHeight,
			)],
			"\n",
		)
		v.Content = components.PanelSpacing.Render(components.RenderPane(
			style.Width(panelWidth).Height(panelHeight),
			"Dropbox sign-in URL",
			body,
			"esc - back  j/k - scroll",
		))
		return v
	}

	question := "Preferred weight unit?"
	choices := "1. kg\n2. lb"
	hint := "Enter to skip (default: kg)."
	title := "Setup (1/2)"
	if m.weightUnit != shared.WeightUnitKg {
		hint = "Enter to skip (keep: " + string(m.weightUnit) + ")."
	}

	if q.step == 1 {
		question = "Link an identity account?"
		choices = "1. Dropbox"
		hint = "Enter to skip (no account saved)."
		title = "Setup (2/2)"
	} else if q.step == 2 {
		question = "Paste the Dropbox code:"
		choices = "ctrl+r retry  ctrl+o URL"
		hint = "Enter to skip (no account saved)."
		title = "Link Dropbox"
	}

	gap := "\n"

	bodyStyle := lipgloss.NewStyle().Width(contentWidth).MaxWidth(contentWidth)
	prefix := bodyStyle.Render(question + gap + choices + gap + hint)
	input := q.input
	input.SetWidth(max(
		1,
		contentWidth-lipgloss.Width(input.Prompt)-1,
	))
	body := prefix + gap + input.View()
	if q.status != "" {
		body += gap + q.status
	}

	if q.authorization != nil {
		body += gap + "Sign-in URL: " + q.authorization.URL
	}

	body = bodyStyle.MaxHeight(contentHeight).Render(body)
	v.Content = components.PanelSpacing.Render(components.RenderPane(
		style.Width(panelWidth).Height(panelHeight),
		title,
		body,
		"esc - cancel",
	))

	inputRow := lipgloss.Height(prefix) - 1 + strings.Count(
		gap,
		"\n",
	)
	if !q.busy && inputRow < contentHeight {
		v.Cursor = input.Cursor()
		if v.Cursor != nil {
			if q.step != 2 {
				value := []rune(input.Value())
				v.Cursor.Position.X = min(
					contentWidth-1,
					lipgloss.Width(input.Prompt)+ansi.StringWidth(string(value[:input.Position()])),
				)
			}
			v.Cursor.Position.X += components.PanelSpacing.GetPaddingLeft() + style.GetBorderLeftSize() + style.GetPaddingLeft()
			v.Cursor.Position.Y = inputRow + components.PanelSpacing.GetPaddingTop() + style.GetBorderTopSize() + style.GetPaddingTop()
		}
	}

	return v
}

func (m RecordsModel) questionnaireURLContent(width int) string {
	return lipgloss.NewStyle().Width(width).MaxWidth(width).Render(
		"Open this URL in your browser:\n" + m.questionnaire.authorization.URL,
	)
}

// END_BLOCK_QUESTIONNAIRE_VIEW

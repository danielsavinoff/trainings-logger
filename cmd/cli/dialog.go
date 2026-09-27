package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

type recordField struct {
	label string
	input textinput.Model
}

type recordDialog struct {
	title     string
	panel     panel
	editing   bool
	inserting bool
	exercise  int
	session   int
	set       int
	fields    []recordField
	active    int
	err       string
}

func newRecordField(label, value, placeholder string) recordField {
	input := textinput.New()
	input.SetVirtualCursor(false)
	input.Prompt = "› "
	input.Placeholder = placeholder
	input.CharLimit = 100
	input.SetValue(value)
	styles := input.Styles()
	styles.Focused.Prompt = lipgloss.NewStyle().Foreground(activePanelColor)
	styles.Focused.Text = lipgloss.NewStyle().Foreground(lipgloss.Color("#F2F4F8"))
	white := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	styles.Blurred.Text = white
	styles.Blurred.Prompt = white
	styles.Blurred.Placeholder = white
	styles.Cursor.Color = activePanelColor
	input.SetStyles(styles)
	return recordField{label: label, input: input}
}

func (m *model) openDialog(editing bool) tea.Cmd {
	exercise := m.selectedExercise()
	if (editing || m.activePanel == setsPanel) && exercise < 0 {
		return nil
	}
	d := &recordDialog{panel: m.activePanel, editing: editing, exercise: exercise}
	action := "New"
	if editing {
		action = "Edit"
	}
	if d.panel == exercisesPanel {
		d.title = action + " exercise"
		name, bodyWeight := "", ""
		if exercise >= 0 {
			e := m.exercises[exercise]
			if editing {
				name = e.name
			}
			if e.bodyWeightKG > 0 {
				bodyWeight = strconv.FormatFloat(e.bodyWeightKG, 'f', -1, 64)
			}
		}
		d.fields = []recordField{
			newRecordField("Name", name, "Exercise name"),
			newRecordField("Body weight (kg)", bodyWeight, "Optional"),
		}
	} else {
		d.title = action + " set"
		set := workoutSet{reps: 8}
		session, index, ok := m.selectedSet()
		if editing && !ok {
			return nil
		}
		if ok {
			d.session, d.set = session, index
			set = m.exercises[exercise].sessions[session].sets[index]
		}
		d.fields = []recordField{
			newRecordField("Reps", strconv.Itoa(set.reps), "8"),
			newRecordField("Weight (kg)", strconv.FormatFloat(set.weightKG, 'f', -1, 64), "0"),
			newRecordField("Rest (m:ss)", fmt.Sprintf("%d:%02d", int(set.rest/time.Second)/60, int(set.rest/time.Second)%60), "2:30"),
		}
	}
	m.goPending = false
	m.dialog = d
	m.resizeDialogInputs()
	return d.focus(0)
}

func (d *recordDialog) focus(index int) tea.Cmd {
	for i := range d.fields {
		d.fields[i].input.Blur()
	}
	d.active = (index + len(d.fields)) % len(d.fields)
	if d.inserting {
		return d.fields[d.active].input.Focus()
	}
	return nil
}

func (m model) updateDialog(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(tea.MouseMsg); ok {
		return m, nil
	}
	if press, ok := msg.(tea.KeyPressMsg); ok {
		switch press.String() {
		case "esc":
			if m.dialog.inserting {
				m.dialog.inserting = false
				m.dialog.fields[m.dialog.active].input.Blur()
			} else {
				m.dialog = nil
			}
			return m, nil
		case "enter":
			cmd := m.saveDialog()
			return m, cmd
		}
		if !m.dialog.inserting {
			switch press.String() {
			case "j", "down", "tab":
				return m, m.dialog.focus(m.dialog.active + 1)
			case "k", "up", "shift+tab":
				return m, m.dialog.focus(m.dialog.active - 1)
			case "i":
				m.dialog.inserting = true
				return m, m.dialog.focus(m.dialog.active)
			}
			return m, nil
		}
		m.dialog.err = ""
	}
	if !m.dialog.inserting {
		return m, nil
	}
	field := &m.dialog.fields[m.dialog.active]
	var cmd tea.Cmd
	field.input, cmd = field.input.Update(msg)
	return m, cmd
}

func (d *recordDialog) invalid(field int, message string) tea.Cmd {
	d.err = message
	return d.focus(field)
}

func nonnegativeNumber(value string) (float64, bool) {
	n, err := strconv.ParseFloat(value, 64)
	return n, err == nil && !math.IsNaN(n) && !math.IsInf(n, 0) && n >= 0
}

func (m *model) saveDialog() tea.Cmd {
	d := m.dialog
	value := func(i int) string { return strings.TrimSpace(d.fields[i].input.Value()) }
	if d.panel == exercisesPanel {
		name := value(0)
		if name == "" {
			return d.invalid(0, "Enter an exercise name.")
		}
		bodyWeight := 0.0
		if value(1) != "" {
			var ok bool
			bodyWeight, ok = nonnegativeNumber(value(1))
			if !ok || bodyWeight == 0 {
				return d.invalid(1, "Body weight must be positive, or left blank.")
			}
		}
		if d.editing {
			m.exercises[d.exercise].name = name
			m.exercises[d.exercise].bodyWeightKG = bodyWeight
		} else {
			m.exercises = append(m.exercises, mockExercise{name: name, bodyWeightKG: bodyWeight})
			d.exercise = len(m.exercises) - 1
			m.setCursor = 0
		}
	} else {
		reps, err := strconv.Atoi(value(0))
		if err != nil || reps < 1 || reps > 255 {
			return d.invalid(0, "Reps must be a whole number from 1 to 255.")
		}
		weight, ok := nonnegativeNumber(value(1))
		if !ok {
			return d.invalid(1, "Weight must be a number of 0 or more.")
		}
		parts := strings.Split(value(2), ":")
		var rest time.Duration
		validRest := false
		if len(parts) == 2 {
			minutes, minuteErr := strconv.ParseUint(parts[0], 10, 32)
			seconds, secondErr := strconv.ParseUint(parts[1], 10, 8)
			if minuteErr == nil && secondErr == nil && seconds < 60 {
				rest, err = time.ParseDuration(fmt.Sprintf("%dm%ds", minutes, seconds))
				validRest = err == nil
			}
		}
		if !validRest {
			return d.invalid(2, "Use minutes:seconds for rest (e.g. 2:30).")
		}
		set := workoutSet{reps: reps, weightKG: weight, rest: rest}
		if d.editing {
			m.exercises[d.exercise].sessions[d.session].sets[d.set] = set
		} else {
			m.appendSet(d.exercise, set, time.Now())
		}
	}
	m.refreshRecords(d.exercise)
	m.dialog = nil
	return nil
}

func (m model) dialogSize() (width int, compact bool) {
	width = max(12, min(64, m.width-4))
	return width, m.height < 16 || width < 44
}

func (m *model) resizeDialogInputs() {
	if m.dialog == nil {
		return
	}
	width, compact := m.dialogSize()
	inputWidth := width - 6 - 20 // Padding, label, prompt and cursor.
	if compact {
		inputWidth = width - 6 - 3
	}
	for i := range m.dialog.fields {
		input := &m.dialog.fields[i].input
		position := input.Position()
		// Reset the horizontal viewport before recalculating it for the new size.
		input.SetWidth(0)
		input.SetCursor(position)
		input.SetWidth(max(1, inputWidth))
		input.CursorEnd()
		input.SetCursor(position)
	}
}

func (m model) dialogView() string {
	d := m.dialog
	width, compact := m.dialogSize()
	innerWidth := max(1, width-6)
	title := d.title
	if compact {
		title += fmt.Sprintf(" · %d/%d", d.active+1, len(d.fields))
	}
	var lines []string
	for i, field := range d.fields {
		if compact && i != d.active {
			continue
		}
		label := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
		if i == d.active {
			label = label.Foreground(activePanelColor)
		}
		if compact {
			input, _ := recordInputView(field, innerWidth)
			lines = append(lines, label.Render(field.label), input)
		} else {
			input, _ := recordInputView(field, innerWidth-17)
			lines = append(lines, label.Width(17).Render(field.label)+input, "")
		}
	}
	if d.err != "" {
		lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color("#F0A3A3")).Width(innerWidth).Render(d.err))
	}
	lines = append(lines, "", m.modalHelpView(innerWidth))
	padding := 1
	if compact {
		padding = 0
	}
	style := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(activePanelColor).
		Foreground(lipgloss.Color("#E6EAF2")).
		Padding(padding, 2).Width(width)
	return renderModalFrame(style, title, strings.Join(lines, "\n"), " Save <enter> ")
}

func (m model) modalHelpView(width int) string {
	bindings := []key.Binding{
		key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "move up")),
		key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "move down")),
		key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "normal mode")),
		key.NewBinding(key.WithKeys("i"), key.WithHelp("i", "insert mode")),
	}
	help := m.table.Help
	help.SetWidth(0)
	center := lipgloss.NewStyle().Width(width).Align(lipgloss.Center)
	var lines []string
	for _, row := range [][]key.Binding{bindings[:2], bindings[2:]} {
		content := help.ShortHelpView(row)
		if lipgloss.Width(content) <= width {
			lines = append(lines, center.Render(content))
			continue
		}
		for _, binding := range row {
			lines = append(lines, center.Render(help.ShortHelpView([]key.Binding{binding})))
		}
	}
	return strings.Join(lines, "\n")
}

// Use cell widths for the visible text and native cursor, including wide Unicode
// characters. Reserving one cell keeps the end-of-input cursor inside the field.
func recordInputView(field recordField, width int) (string, int) {
	input := field.input
	promptWidth := lipgloss.Width(input.Prompt)
	textWidth := max(1, width-promptWidth)
	value := input.Value()
	position := min(input.Position(), len([]rune(value)))
	caret := ansi.StringWidth(string([]rune(value)[:position]))
	start := max(0, caret-textWidth+1)
	styles := input.Styles().Blurred
	if input.Focused() {
		styles = input.Styles().Focused
	}
	textStyle := styles.Text
	if value == "" {
		value = input.Placeholder
		textStyle = styles.Placeholder
	}
	visible := ansi.Truncate(ansi.Cut(value, start, start+textWidth), textWidth, "")
	visible += strings.Repeat(" ", max(0, textWidth-lipgloss.Width(visible)))
	return ansi.Truncate(styles.Prompt.Render(input.Prompt)+textStyle.Render(visible), width, ""), promptWidth + caret - start
}

func renderModalFrame(style lipgloss.Style, title, body, footer string) string {
	rendered := strings.Split(style.Render(body), "\n")
	border := style.GetBorderStyle()
	borderWidth := max(0, lipgloss.Width(rendered[0])-2)
	label := []rune(" " + title + " ")
	if len(label) > borderWidth {
		label = label[:borderWidth]
	}
	remaining := borderWidth - lipgloss.Width(string(label))
	borderStyle := lipgloss.NewStyle().Foreground(style.GetBorderTopForeground())
	rendered[0] = borderStyle.Render(border.TopLeft+strings.Repeat(border.Top, remaining/2)) +
		lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Render(string(label)) +
		borderStyle.Render(strings.Repeat(border.Top, remaining-remaining/2)+border.TopRight)
	footer = footer[:min(len(footer), borderWidth)]
	remaining = borderWidth - len(footer)
	white := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	rendered[len(rendered)-1] = borderStyle.Render(border.BottomLeft+strings.Repeat(border.Bottom, remaining/2)) +
		white.Render(footer) + borderStyle.Render(strings.Repeat(border.Bottom, remaining-remaining/2)+border.BottomRight)
	return strings.Join(rendered, "\n")
}

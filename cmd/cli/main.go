package main

import (
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/paginator"
	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var activePanelColor = lipgloss.Color("#8BA4D4")

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("245")).
	Padding(1, 2)

var panelSpacing = lipgloss.NewStyle().Padding(1, 2)

type panel int

const (
	exercisesPanel panel = iota
	setsPanel
)

type model struct {
	exercises     []mockExercise
	table         table.Model
	timeline      viewport.Model
	exercisePager paginator.Model
	setPager      paginator.Model
	setPageStarts []int
	activePanel   panel
	goPending     bool
	setCursor     int
	setLines      []int
	width         int
	height        int
	dialog        *recordDialog
}

func (m model) Init() tea.Cmd { return nil }

func (m model) resize() model {
	helpHeight := lipgloss.Height(m.helpView())

	frameWidth := baseStyle.GetHorizontalFrameSize() + panelSpacing.GetHorizontalFrameSize()
	width := max(0, m.width/2-frameWidth)
	rightWidth := max(0, m.width-m.width/2-frameWidth)
	panelHeight := max(0, m.height-helpHeight-panelSpacing.GetVerticalFrameSize())
	frameHeight := baseStyle.GetVerticalFrameSize()

	m.table.SetColumns(evenColumns(m.table.Columns(), width))
	m.table.SetWidth(width)
	// Reserve at least one line for the column names.
	m.table.SetHeight(max(1, panelHeight-frameHeight-1))
	m.timeline.SetWidth(rightWidth)
	// Keep the Sets column headings fixed above the scrolling history.
	m.timeline.SetHeight(max(0, panelHeight-frameHeight-2))
	m.syncExercisePager()
	m.refreshTimeline()
	m.resizeDialogInputs()

	return m
}

func evenColumns(source []table.Column, width int) []table.Column {
	columns := append([]table.Column(nil), source...)
	if len(columns) > 0 {
		// Column widths exclude the padding applied to each cell.
		cellPadding := table.DefaultStyles().Cell.GetHorizontalFrameSize()
		for i := range columns {
			columnWidth := width / len(columns)
			if i < width%len(columns) {
				columnWidth++
			}
			columns[i].Width = max(0, columnWidth-cellPadding)
		}
	}
	return columns
}

func (m *model) focusPanel(next panel) {
	m.activePanel = next
	m.table.Blur()
	if next == exercisesPanel {
		m.table.Focus()
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if size, ok := msg.(tea.WindowSizeMsg); ok {
		m.width, m.height = size.Width, size.Height
		m = m.resize()
		return m, nil
	}
	if press, ok := msg.(tea.KeyPressMsg); ok && press.String() == "ctrl+c" {
		return m, tea.Quit
	}
	if m.dialog != nil {
		return m.updateDialog(msg)
	}
	if wheel, ok := msg.(tea.MouseWheelMsg); ok {
		return m.scrollPanel(wheel)
	}
	// Mouse movement and clicks must not scroll the table's internal viewport.
	if _, ok := msg.(tea.MouseMsg); ok {
		return m, nil
	}
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if m.goPending {
			m.goPending = false
			switch msg.String() {
			case "e":
				m.focusPanel(exercisesPanel)
				return m, nil
			case "s":
				m.focusPanel(setsPanel)
				return m, nil
			}
		} else if msg.String() == "g" {
			m.goPending = true
			return m, nil
		}
		switch msg.String() {
		case "esc", "tab", "shift+tab":
			return m, nil
		case "q":
			return m, tea.Quit
		case "n":
			cmd = m.openDialog(false)
			return m, cmd
		case "e":
			cmd = m.openDialog(true)
			return m, cmd
		case "d":
			m.deleteRecord()
			return m, nil
		case "enter":
			if m.activePanel == exercisesPanel {
				m.focusPanel(setsPanel)
			}
			return m, nil
		}
	}
	if press, ok := msg.(tea.KeyPressMsg); ok && m.paginate(press) {
		return m, nil
	}
	if m.activePanel == setsPanel {
		if press, ok := msg.(tea.KeyPressMsg); ok && len(m.setLines) > 0 {
			keys := m.table.KeyMap
			next := m.setCursor
			switch {
			case key.Matches(press, keys.LineUp):
				next--
			case key.Matches(press, keys.LineDown):
				next++
			case key.Matches(press, keys.HalfPageUp):
				next -= max(1, m.timeline.Height()/2)
			case key.Matches(press, keys.HalfPageDown):
				next += max(1, m.timeline.Height()/2)
			case key.Matches(press, keys.GotoTop):
				next = 0
			case key.Matches(press, keys.GotoBottom):
				next = len(m.setLines) - 1
			}
			m.setCursor = min(max(0, next), len(m.setLines)-1)
			m.refreshTimeline()
			return m, nil
		}
		return m, nil
	}
	previous := m.table.Cursor()
	m.table, cmd = m.table.Update(msg)
	m.syncExercisePager()
	if m.table.Cursor() != previous {
		m.setCursor = 0
		m.refreshTimeline()
	}
	return m, cmd
}

func newPaginator() paginator.Model {
	p := paginator.New()
	p.Type = paginator.Dots
	p.ActiveDot = lipgloss.NewStyle().Foreground(activePanelColor).Render("● ")
	p.InactiveDot = lipgloss.NewStyle().Foreground(baseStyle.GetBorderTopForeground()).Render("○ ")
	p.KeyMap.PrevPage = key.NewBinding(key.WithKeys("left", "h", "pgup", "b"))
	p.KeyMap.NextPage = key.NewBinding(key.WithKeys("right", "l", "pgdown", "f", "space"))
	return p
}

func (m *model) syncExercisePager() {
	m.exercisePager.PerPage = max(1, m.table.Height())
	m.exercisePager.SetTotalPages(max(1, len(m.table.Rows())))
	m.exercisePager.Page = max(0, m.table.Cursor()) / m.exercisePager.PerPage
}

func (m *model) paginate(msg tea.KeyPressMsg) bool {
	p := &m.exercisePager
	if m.activePanel == setsPanel {
		p = &m.setPager
	}
	if !key.Matches(msg, p.KeyMap.PrevPage, p.KeyMap.NextPage) {
		return false
	}
	previous := p.Page
	*p, _ = p.Update(msg)
	if p.Page == previous {
		return true
	}
	if m.activePanel == exercisesPanel {
		m.table.SetCursor(p.Page * p.PerPage)
		m.setCursor = 0
	} else {
		m.setCursor = m.setPageStarts[p.Page]
	}
	m.refreshTimeline()
	return true
}

func paginatorView(p paginator.Model, width int) string {
	// Long histories use numbers so the indicator always fits on one line.
	if p.TotalPages*2 > width {
		p.Type = paginator.Arabic
	}
	return lipgloss.NewStyle().Width(width).MaxWidth(width).Align(lipgloss.Center).Render(p.View())
}

func (m model) exercisesView() string {
	start, end := m.exercisePager.GetSliceBounds(len(m.table.Rows()))
	t := m.table
	t.SetRows(m.table.Rows()[start:end])
	t.GotoTop()
	t.SetCursor(max(0, m.table.Cursor()-start))
	return lipgloss.JoinVertical(lipgloss.Left, t.View(), paginatorView(m.exercisePager, m.table.Width()))
}

func (m model) helpView() string {
	bindings := []key.Binding{
		key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "move up")),
		key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "move down")),
		key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "move left")),
		key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "move right")),
	}
	actions := []key.Binding{
		key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new record")),
		key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit record")),
		key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete record")),
		key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	}
	help := m.table.Help
	help.SetWidth(m.width)
	center := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center)
	return lipgloss.JoinVertical(lipgloss.Left,
		center.Render(help.ShortHelpView(bindings)),
		center.Render(help.ShortHelpView(actions)),
	)
}

func (m model) View() tea.View {
	styles := []lipgloss.Style{baseStyle, baseStyle}
	styles[m.activePanel] = styles[m.activePanel].BorderForeground(activePanelColor)
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top,
			panelSpacing.Render(renderPane(styles[exercisesPanel], "Exercises <ge>", m.exercisesView(), position(m.table.Cursor(), len(m.table.Rows())))),
			panelSpacing.Render(renderPane(styles[setsPanel], "Sets <gs>", m.setsView(), position(m.setCursor, len(m.setLines)))),
		),
		m.helpView(),
	)
	if m.width > 0 && m.height > 0 {
		content = lipgloss.NewStyle().MaxWidth(m.width).MaxHeight(m.height).Render(content)
	}

	layers := []*lipgloss.Layer{lipgloss.NewLayer(content).ID("panels")}
	var cursor *tea.Cursor
	if m.dialog != nil {
		var dialogLayers []*lipgloss.Layer
		dialogLayers, cursor = m.dialogLayers()
		layers = append(layers, dialogLayers...)
	}
	canvas := lipgloss.NewCanvas(max(0, m.width), max(0, m.height))
	content = canvas.Compose(lipgloss.NewCompositor(layers...)).Render()
	v := tea.NewView(content)
	v.Cursor = cursor
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion

	return v
}

func position(cursor, count int) string {
	if count == 0 {
		return "0 of 0"
	}
	return fmt.Sprintf("%d of %d", cursor+1, count)
}

func renderPane(style lipgloss.Style, title, content, counter string) string {
	rendered := style.Render(content)
	lines := strings.Split(rendered, "\n")
	border := style.GetBorderStyle()
	width := max(0, lipgloss.Width(lines[0])-2)
	label := []rune(" " + title + " ")
	if len(label) > width {
		label = label[:width]
	}
	topStyle := lipgloss.NewStyle().Foreground(style.GetBorderTopForeground())
	lines[0] = topStyle.Render(border.TopLeft) + string(label) +
		topStyle.Render(strings.Repeat(border.Top, width-len(label))+border.TopRight)
	footer := []rune(" " + counter + " ")
	if len(footer) > width {
		footer = footer[:width]
	}
	bottom := border.BottomLeft + strings.Repeat(border.Bottom, width-len(footer)) + string(footer) + border.BottomRight
	lines[len(lines)-1] = lipgloss.NewStyle().Foreground(style.GetBorderBottomForeground()).Render(bottom)
	return strings.Join(lines, "\n")
}

func dateLine(date string, width int) string {
	date = " " + date + " "
	if width <= len(date) {
		return date[:max(0, width)]
	}
	remaining := width - len(date)
	return strings.Repeat("─", remaining/2) + date + strings.Repeat("─", remaining-remaining/2)
}

func (m model) setColumns() []table.Column {
	return evenColumns([]table.Column{{Title: "Reps"}, {Title: "Exercise Weight"}, {Title: "Rest (Time)"}}, m.timeline.Width())
}

func (m model) renderSetRow(values []string, style lipgloss.Style) string {
	columns := m.setColumns()
	cells := make([]string, 0, len(columns))
	for i, column := range columns {
		if column.Width <= 0 {
			continue
		}
		cell := lipgloss.NewStyle().Width(column.Width).MaxWidth(column.Width).Inline(true).Render(values[i])
		cells = append(cells, lipgloss.NewStyle().Padding(0, 1).Render(cell))
	}
	return style.Render(lipgloss.JoinHorizontal(lipgloss.Top, cells...))
}

func (m model) setsView() string {
	columns := m.setColumns()
	titles := make([]string, len(columns))
	for i, column := range columns {
		titles[i] = column.Title
	}
	header := m.renderSetRow(titles, lipgloss.NewStyle().Foreground(activePanelColor))
	return lipgloss.JoinVertical(lipgloss.Left, header, m.timeline.View(), paginatorView(m.setPager, m.timeline.Width()))
}

func (m *model) refreshTimeline() {
	m.setLines = nil
	m.setPageStarts = nil
	m.setPager.Page = 0
	m.setPager.TotalPages = 1
	m.setPager.PerPage = max(1, m.timeline.Height())
	if len(m.exercises) == 0 {
		m.setCursor = 0
		m.timeline.SetContent("")
		return
	}
	exercise := m.exercises[m.table.Cursor()]
	count := 0
	for _, session := range exercise.sessions {
		count += len(session.sets)
	}
	m.setCursor = max(0, min(m.setCursor, count-1))
	lines := []string{}
	index := 0
	for _, session := range exercise.sessions {
		separator := dateLine(session.date.Format("2 January, 2006"), m.timeline.Width())
		for sessionIndex, set := range session.sets {
			pageHeight := m.setPager.PerPage
			needsDate := sessionIndex == 0 || len(lines)%pageHeight == 0
			needed := 1
			if needsDate && pageHeight > 1 {
				needed++
			}
			if remaining := pageHeight - len(lines)%pageHeight; remaining < needed {
				lines = append(lines, make([]string, remaining)...)
				needsDate = true
			}
			if len(lines)%pageHeight == 0 {
				m.setPageStarts = append(m.setPageStarts, index)
			}
			if needsDate && pageHeight > 1 {
				lines = append(lines, lipgloss.NewStyle().Foreground(baseStyle.GetBorderTopForeground()).Render(separator))
			}
			m.setLines = append(m.setLines, len(lines))
			style := lipgloss.NewStyle()
			if index == m.setCursor {
				style = selectedStyle()
			}
			restSeconds := int(set.rest.Seconds())
			values := []string{
				fmt.Sprint(set.reps),
				fmt.Sprintf("%g KG", set.weightKG),
				fmt.Sprintf("%d:%02d", restSeconds/60, restSeconds%60),
			}
			lines = append(lines, m.renderSetRow(values, style))
			index++
		}
	}
	if len(m.setLines) > 0 {
		m.setCursor = min(max(0, m.setCursor), len(m.setLines)-1)
		m.setPager.SetTotalPages(len(lines))
		m.setPager.Page = m.setLines[m.setCursor] / m.setPager.PerPage
		// Pad the last page so the viewport can align it at the top.
		lines = append(lines, make([]string, m.setPager.TotalPages*m.setPager.PerPage-len(lines))...)
	}
	m.timeline.SetContent(strings.Join(lines, "\n"))
	m.timeline.SetYOffset(m.setPager.Page * m.setPager.PerPage)
}

func main() {
	m := newModel(mockExercises())
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func selectedStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#182235")).Background(activePanelColor)
}

func newModel(exercises []mockExercise) model {
	columns := []table.Column{{Title: "Name"}, {Title: "Body Weight %"}}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(exerciseRows(exercises)),
		table.WithFocused(true),
		table.WithHeight(7),
		table.WithWidth(42),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		Foreground(activePanelColor).
		Bold(false)
	s.Selected = selectedStyle()
	t.SetStyles(s)

	m := model{exercises: exercises, table: t, timeline: viewport.New(), exercisePager: newPaginator(), setPager: newPaginator(), width: 100, height: 24}
	return m.resize()
}

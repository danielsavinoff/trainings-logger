package components

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

type Field struct {
	label string
	input textinput.Model
}

func NewField(label, value, placeholder string) Field {
	input := textinput.New()
	input.SetVirtualCursor(false)
	input.Prompt = "› "
	input.Placeholder = placeholder
	input.CharLimit = 100
	input.SetValue(value)

	styles := input.Styles()
	styles.Focused.Prompt = lipgloss.NewStyle().Foreground(ActivePanelColor)
	styles.Focused.Text = lipgloss.NewStyle().Foreground(lipgloss.Color("#F2F4F8"))

	white := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	styles.Blurred.Text = white
	styles.Blurred.Prompt = white
	styles.Blurred.Placeholder = white
	styles.Cursor.Color = ActivePanelColor
	input.SetStyles(styles)

	return Field{label: label, input: input}
}

// Use cell widths for the visible text and native cursor, including wide Unicode
// characters. Reserving one cell keeps the end-of-input cursor inside the field.
func fieldView(field Field, width int) (string, int) {
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
	if value == "" && input.Focused() {
		value = input.Placeholder
		textStyle = styles.Placeholder
	}

	visible := ansi.Truncate(ansi.Cut(value, start, start+textWidth), textWidth, "")
	visible += strings.Repeat(" ", max(0, textWidth-lipgloss.Width(visible)))

	return ansi.Truncate(styles.Prompt.Render(input.Prompt)+textStyle.Render(visible), width, ""), promptWidth + caret - start
}

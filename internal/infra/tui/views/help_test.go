package views

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

func TestHelpRendering(t *testing.T) {
	for _, size := range [][2]int{{100, 40}, {161, 40}, {40, 12}, {12, 8}, {1, 1}, {0, 0}} {
		m := newHelpModel(
			size[0],
			size[1],
		)
		view := ansi.Strip(m.view())
		if lipgloss.Width(view) > size[0] || (view != "" && lipgloss.Height(view) > size[1]) {
			t.Fatalf(
				"help exceeds %v: %q",
				size,
				view,
			)
		}
		if strings.Contains(
			view,
			"esc/?",
		) || strings.Contains(
			view,
			"q  quit",
		) || strings.HasPrefix(
			view,
			"Help\n",
		) {
			t.Fatalf(
				"unexpected page title or footer: %q",
				view,
			)
		}
	}
	wide := ansi.Strip(newHelpModel(
		161,
		40,
	).view())
	if !strings.Contains(
		wide,
		"Navigation",
	) || !strings.Contains(
		wide,
		"Pages",
	) || !strings.Contains(
		wide,
		"Records (active panel)",
	) || !strings.Contains(
		wide,
		"Help",
	) {
		t.Fatalf(
			"missing group headings: %q",
			wide,
		)
	}
	if !strings.Contains(
		wide,
		"ge       focus exercises",
	) || !strings.Contains(
		wide,
		"n    new record",
	) {
		t.Fatalf(
			"keys and explanations are not tab aligned: %q",
			wide,
		)
	}
	if !strings.Contains(
		wide,
		"Navigation",
	) || !strings.Contains(
		wide,
		"Pages",
	) {
		t.Fatal("wide view should put groups side by side")
	}
	compact := ansi.Strip(newHelpModel(
		40,
		12,
	).view())
	if strings.Index(
		compact,
		"Pages",
	) < strings.Index(
		compact,
		"focus sets",
	) {
		t.Fatal("compact view should stack groups")
	}
	colored := newHelpModel(
		161,
		40,
	).view()
	if !strings.Contains(
		colored,
		"38;2;139;164;212mge",
	) {
		t.Fatal("shortcut keys should use active panel color")
	}
}

func TestHelpScrollingAndResize(t *testing.T) {
	m := newHelpModel(
		40,
		8,
	)
	for i := 0; i < 100; i++ {
		m = m.update(tea.MouseWheelMsg{Button: tea.MouseWheelDown})
	}
	if !m.viewport.AtBottom() {
		t.Fatal("mouse wheel did not reach the bottom")
	}
	m = m.update(tea.MouseWheelMsg{Button: tea.MouseWheelUp})
	if m.viewport.AtBottom() {
		t.Fatal("mouse wheel up did not scroll")
	}
	m.resize(
		100,
		40,
	)
	if !strings.Contains(
		ansi.Strip(m.view()),
		"Navigation",
	) {
		t.Fatal("resize did not expose help groups")
	}
	m.resize(
		0,
		0,
	)
	if m.view() != "" {
		t.Fatal("zero dimensions must render empty")
	}
}

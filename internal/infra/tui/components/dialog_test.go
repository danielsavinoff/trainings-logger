package components

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/danielsavinoff/trainings-logger/internal/infra/tui/testutils"
)

func TestDialogModesValidationAndResize(t *testing.T) {
	d := NewDialog("Exercise", []Field{NewField("Name", "Squat", ""), NewField("Share", "75", "")})
	d.Resize(100, 24)
	d.Init()
	press := func(keys ...string) {
		for _, key := range keys {
			d.Update(testutils.Key(key))
		}
	}

	press("j", "i", "ctrl+u", "101", "enter")
	if d.State != DialogSubmitted || d.Value(1) != "101" {
		t.Fatal("insert mode did not submit the edited field")
	}

	d.Invalid(1, "Enter 0–100.")
	if d.State != DialogOpen || !strings.Contains(d.View(), "Enter 0–100.") {
		t.Fatal("validation did not reopen the form with an error")
	}

	d.Resize(40, 12)
	_, cursor := d.Layers()
	if cursor == nil || cursor.Shape != tea.CursorBar || cursor.Position.X >= 40 || cursor.Position.Y >= 12 {
		t.Fatal("resizing lost the insert cursor")
	}

	press("ctrl+u", "50")
	if strings.Contains(d.View(), "Enter 0–100.") {
		t.Fatal("editing did not clear the validation error")
	}

	press("esc")
	_, cursor = d.Layers()
	if d.State != DialogOpen || cursor.Shape != tea.CursorBlock {
		t.Fatal("first escape should return to normal mode")
	}

	press("esc")
	if d.State != DialogClosed {
		t.Fatal("second escape should close the form")
	}
}

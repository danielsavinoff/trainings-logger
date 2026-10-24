// START_MODULE_CONTRACT
//   PURPOSE: Keep modal sizing and native cursors aligned with shared compact padding.
//   SCOPE: Terminal bounds, field sizing, centered layers and input cursor offsets.
//   DEPENDS: none
//   LINKS: M-TUI-COMPONENTS, V-M-TUI-COMPONENTS
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   dialogSize - choose bounded modal width and compact mode
//   Resize - update field widths
//   Layers - center the bounded modal and cursor
//   dialogCursor - locate the active field using shared padding
// END_MODULE_MAP

package components

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// START_BLOCK_DIALOG_LAYOUT

func (d *Dialog) dialogSize() (width int, compact bool) {
	width = max(
		0,
		min(
			64,
			d.width-PanelSpacing.GetHorizontalFrameSize(),
		),
	)

	return width, d.height < 16 || width < 44
}

func (d *Dialog) Resize(width, height int) {
	d.width, d.height = width, height

	width, compact := d.dialogSize()
	innerWidth := width - BaseStyle.GetHorizontalFrameSize()
	inputWidth := innerWidth - 17 - 3 // Label, prompt and cursor.
	if compact {
		inputWidth = innerWidth - 3
	}

	for i := range d.fields {
		input := &d.fields[i].input
		position := input.Position()

		input.SetWidth(0)
		input.SetCursor(position)
		input.SetWidth(max(
			1,
			inputWidth,
		))
		input.CursorEnd()
		input.SetCursor(position)
	}
}

func (d *Dialog) Layers() ([]*lipgloss.Layer, *tea.Cursor) {
	form := d.View()
	if form == "" {
		return nil, nil
	}

	formWidth, formHeight := lipgloss.Size(form)
	x, y := max(
		0,
		(d.width-formWidth)/2,
	), max(
		0,
		(d.height-formHeight)/2,
	)
	layers := []*lipgloss.Layer{lipgloss.NewLayer(form).ID("dialog").Z(1).X(x).Y(y)}

	cursor := d.dialogCursor(
		x,
		y,
	)
	if cursor != nil && cursor.Position.Y >= y+formHeight-1 {
		cursor = nil
	}

	return layers, cursor
}

func (d *Dialog) dialogCursor(x, y int) *tea.Cursor {
	width, compact := d.dialogSize()
	innerWidth := width - BaseStyle.GetHorizontalFrameSize()
	contentX := BaseStyle.GetBorderLeftSize() + BaseStyle.GetPaddingLeft()
	contentY := BaseStyle.GetBorderTopSize() + BaseStyle.GetPaddingTop()
	fieldWidth := innerWidth - 17
	fieldX, fieldY := contentX+17, contentY+d.active
	if compact {
		fieldWidth = innerWidth
		fieldX, fieldY = contentX, contentY+1
	}

	_, caret := fieldView(
		d.fields[d.active],
		fieldWidth,
	)
	cursor := tea.NewCursor(
		x+fieldX+caret,
		y+fieldY,
	)
	cursor.Color = ActivePanelColor
	cursor.Shape = tea.CursorBlock
	cursor.Blink = false
	if d.inserting {
		cursor.Shape = tea.CursorBar
		cursor.Blink = true
	}

	if cursor.Position.X >= d.width || cursor.Position.Y >= d.height {
		return nil
	}

	return cursor
}

// END_BLOCK_DIALOG_LAYOUT

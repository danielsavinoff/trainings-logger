// START_MODULE_CONTRACT
//   PURPOSE: Define records action bindings.
//   SCOPE: New, edit, delete, help, and local quit shortcuts.
//   DEPENDS: none
//   LINKS: M-TUI-RECORDS, V-M-TUI-RECORDS
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   recordsKeyMap - hold action bindings
//   defaultRecordsKeyMap - construct records shortcuts
// END_MODULE_MAP

package views

import "charm.land/bubbles/v2/key"

// START_BLOCK_RECORDS_KEYS
type recordsKeyMap struct {
	New, Edit, Delete, Help, Quit key.Binding
}

func defaultRecordsKeyMap() recordsKeyMap {
	return recordsKeyMap{
		New: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp(
				"n",
				"new record",
			),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp(
				"e",
				"edit record",
			),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp(
				"d",
				"delete record",
			),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp(
				"?",
				"help",
			),
		),
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp(
				"q",
				"quit",
			),
		),
	}
}

// END_BLOCK_RECORDS_KEYS

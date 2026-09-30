package views

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

func (m *RecordsModel) paginate(msg tea.KeyPressMsg) bool {
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
		m.refreshSelectedSets()
	} else {
		m.setCursor = m.setPageStarts[p.Page]
		m.refreshTimeline()
	}

	return true
}

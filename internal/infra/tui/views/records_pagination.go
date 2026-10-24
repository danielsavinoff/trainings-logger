// START_MODULE_CONTRACT
//   PURPOSE: Navigate pages in the active records table.
//   SCOPE: Four independent page states without separate paginator indicators.
//   DEPENDS: M-APP-SESSIONS, M-APP-PERFORMANCES, M-DOMAIN-SESSION, M-DOMAIN-PERFORMANCE, M-DOMAIN-EXERCISE, M-DOMAIN-SET
//   LINKS: M-TUI-RECORDS, V-M-TUI-RECORDS
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   activeTable - select the focused table
//   activePager - select the focused pager
//   syncPager - derive paging state for one table
//   syncPagers - synchronize all four pagers
//   pageTableView - render a single table page
//   paginate - move the active page
//   selectionChanged - reload children only when their selected parent changes
// END_MODULE_MAP

package views

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/paginator"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
)

// START_BLOCK_RECORDS_PAGINATION
func (m *RecordsModel) activeTable() *table.Model {
	switch m.activePanel {
	case catalogPanel:
		return &m.catalogTable
	case performancesPanel:
		return &m.performanceTable
	case setsPanel:
		return &m.setTable
	default:
		return &m.sessionTable
	}
}

func (m *RecordsModel) activePager() *paginator.Model {
	switch m.activePanel {
	case catalogPanel:
		return &m.catalogPager
	case performancesPanel:
		return &m.performancePager
	case setsPanel:
		return &m.setPager
	default:
		return &m.sessionPager
	}
}

func syncPager(p *paginator.Model, t table.Model) {
	p.PerPage = max(
		1,
		t.Height(),
	)
	p.SetTotalPages(max(
		1,
		len(t.Rows()),
	))
	p.Page = max(
		0,
		t.Cursor(),
	) / p.PerPage
}

func (m *RecordsModel) syncPagers() {
	syncPager(
		&m.catalogPager,
		m.catalogTable,
	)
	syncPager(
		&m.performancePager,
		m.performanceTable,
	)
	syncPager(
		&m.sessionPager,
		m.sessionTable,
	)
	syncPager(
		&m.setPager,
		m.setTable,
	)
}

func pageTableView(t table.Model, p paginator.Model) string {
	start, end := p.GetSliceBounds(len(t.Rows()))
	cursor := t.Cursor()
	t.SetRows(t.Rows()[start:end])
	t.GotoTop()
	t.SetCursor(max(
		0,
		cursor-start,
	))
	return t.View()
}

func (m *RecordsModel) paginate(msg tea.KeyPressMsg) bool {
	p := m.activePager()
	if !key.Matches(
		msg,
		p.KeyMap.PrevPage,
		p.KeyMap.NextPage,
	) {
		return false
	}

	previous := p.Page
	*p, _ = p.Update(msg)
	if p.Page != previous {
		m.activeTable().SetCursor(p.Page * p.PerPage)
		m.selectionChanged()
	}

	return true
}

func (m *RecordsModel) selectionChanged() {
	m.syncPagers()
	switch m.activePanel {
	case sessionsPanel:
		m.refreshSelectedPerformances()
	case performancesPanel:
		m.refreshSelectedSets()
	}

	*m = m.resize()
}

// END_BLOCK_RECORDS_PAGINATION

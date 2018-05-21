package main

import (
	"fmt"

	ui "github.com/dpetzold/termui"
)

func NewEventsPanel(height int) *ui.Table {
	p := ui.NewTable()
	p.Height = height
	p.BorderLabel = "Events"
	p.TextAlign = ui.AlignLeft
	p.Separator = false
	p.Headers = true
	p.Analysis()
	return p
}

func EventsFooter() *ui.Par {
	return DashboardFooter()
}

func ShowEvents() {
	Globals.EventsPanel.Height = ui.TermHeight() - 1
	ui.Body.Rows = []*ui.Row{
		ui.NewRow(ui.NewCol(12, 0, Globals.EventsPanel)),
		ui.NewRow(ui.NewCol(12, 0, EventsFooter())),
	}
	Globals.ActiveWindow = "EventsWindow"
}

func updateEvents(eventsPanel *ui.Table) {
	eventRows := [][]string{
		[]string{"Last Seen", "Count", "Name", "Kind", "Type", "Reason", "Message"},
	}

	events := Globals.Events

	for _, e := range events {
		eventRows = append(eventRows, []string{
			timeToDurationStr(e.LastTimestamp.Time),
			fmt.Sprintf("%d", e.Count),
			e.ObjectMeta.Name[0:20],
			e.InvolvedObject.Kind,
			e.Type,
			e.Reason,
			e.Message,
		})
	}

	max_rows := EVENTS_PANEL_HEIGHT - 3
	if len(eventRows) > max_rows {
		eventRows = eventRows[0:max_rows]
	}

	eventsPanel.Rows = eventRows
	eventsPanel.Analysis()
}

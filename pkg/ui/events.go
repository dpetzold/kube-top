package ui

import (
	"fmt"

	"github.com/dpetzold/kube-top/pkg/global"
	"github.com/dpetzold/termui"
)

func NewEventsPanel(height int) *termui.Table {
	p := termui.NewTable()
	p.Height = height
	p.BorderLabel = "Events"
	p.TextAlign = termui.AlignLeft
	p.Separator = false
	p.Headers = true
	p.Analysis()
	return p
}

func EventsFooter() *termui.Par {
	return DashboardFooter()
}

func ShowEvents() {
	global.EventsPanel.Height = termui.TermHeight() - 1
	termui.Body.Rows = []*termui.Row{
		termui.NewRow(termui.NewCol(12, 0, global.EventsPanel)),
		termui.NewRow(termui.NewCol(12, 0, EventsFooter())),
	}
	global.ActiveWindow = "EventsWindow"
}

func updateEvents(eventsPanel *termui.Table) {
	eventRows := [][]string{
		[]string{"Last Seen", "Count", "Name", "Kind", "Type", "Reason", "Message"},
	}

	events := global.Events

	for _, e := range events {
		eventRows = append(eventRows, []string{
			TimeToDurationStr(e.LastTimestamp.Time),
			fmt.Sprintf("%d", e.Count),
			e.ObjectMeta.Name[0:20],
			e.InvolvedObject.Kind,
			e.Type,
			e.Reason,
			e.Message,
		})
	}

	max_rows := global.EventsPanel.Height - 3
	if len(eventRows) > max_rows {
		eventRows = eventRows[0:max_rows]
	}

	eventsPanel.Rows = eventRows
	eventsPanel.Analysis()
}
